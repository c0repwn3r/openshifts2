package routes

import (
	"github.com/c0repwn3r/openshifts/api/apierr"
	"github.com/c0repwn3r/openshifts/api/email"
	"github.com/c0repwn3r/openshifts/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`

	OrgName string `json:"org_name"`
}

type RegisterResponse struct {
	Organization models.Organization `json:"organization"`
	User         models.User         `json:"user"`
	Token        string              `json:"token"`
}

func (e *Env) RegisterRoute(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.BindJSON(&registerRequest); err != nil {
		return
	}

	var maybeExistingUser models.User
	result := e.DB.Where("email = ?", registerRequest.Email).First(&maybeExistingUser)
	if result.RowsAffected == 1 {
		c.JSON(http.StatusConflict, apierr.ErrorResponse{
			Errors: []apierr.Error{
				{
					Code:    apierr.EEmailInUse,
					Message: "user with this email already exists",
					Path:    "email",
				},
			},
		})
		return
	}

	if email.IsFreeEmailProvider(email.GetDomain(registerRequest.Email)) {
		c.JSON(http.StatusBadRequest, apierr.ErrorResponse{
			Errors: []apierr.Error{
				{
					Code:    apierr.EFreeEmail,
					Message: "email is a free email provider",
					Path:    "email",
				},
			},
		})
		return
	}

	var org models.Organization
	result = e.DB.Where("email_domain = ?", email.GetDomain(registerRequest.Email)).First(&org)
	if result.RowsAffected == 1 {
		c.JSON(http.StatusConflict, apierr.ErrorResponse{
			Errors: []apierr.Error{
				{
					Code:    apierr.EDomainInUse,
					Message: "email domain in use by another organization",
					Path:    "email",
				},
			},
		})
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierr.ErrorResponse{
			Errors: []apierr.Error{
				{
					Code:    apierr.EPasswordTooLong,
					Message: "password is too long",
					Path:    "password",
				},
			},
		})
		return
	}

	newOrg := &models.Organization{
		EmailDomain: email.GetDomain(registerRequest.Email),
		Name:        registerRequest.OrgName,
	}
	e.DB.Create(&newOrg)

	hours := models.OrganizationHours{
		Days: models.WeekdayAvailability{
			Sunday: []models.Availability{
				{Start: models.HMToMinute(9, 0), End: models.HMToMinute(17, 0)},
			},
			Monday: []models.Availability{
				{Start: models.HMToMinute(9, 0), End: models.HMToMinute(17, 0)},
			},
			Tuesday: []models.Availability{
				{Start: models.HMToMinute(9, 0), End: models.HMToMinute(17, 0)},
			},
			Wednesday: []models.Availability{
				{Start: models.HMToMinute(9, 0), End: models.HMToMinute(17, 0)},
			},
			Thursday: []models.Availability{
				{Start: models.HMToMinute(9, 0), End: models.HMToMinute(17, 0)},
			},
			Friday: []models.Availability{
				{Start: models.HMToMinute(9, 0), End: models.HMToMinute(17, 0)},
			},
			Saturday: []models.Availability{
				{Start: models.HMToMinute(9, 0), End: models.HMToMinute(17, 0)},
			},
		},
		Overrides: []models.DayOverride{},
	}
	models.SetOrgHours(&e.DB, newOrg.ID, hours)

	newUser := &models.User{
		FirstName:       registerRequest.FirstName,
		LastName:        registerRequest.LastName,
		Email:           registerRequest.Email,
		PermissionLevel: models.PermissionLevelAdmin,
		Password:        string(h),
		OrganizationID:  newOrg.ID,
	}
	e.DB.Create(&newUser)

	newToken := &models.Token{
		ID:     uuid.New().String(),
		UserID: newUser.ID,
	}
	e.DB.Create(&newToken)

	c.JSON(http.StatusOK, RegisterResponse{
		Organization: *newOrg,
		User:         *newUser,
		Token:        newToken.ID,
	})
}
