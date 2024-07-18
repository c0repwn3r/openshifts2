package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type PermissionLevel int64

const (
	PermissionLevelEmployee PermissionLevel = iota
	PermissionLevelManager
	PermissionLevelAdmin
)

func (p PermissionLevel) String() string {
	switch p {
	case PermissionLevelEmployee:
		return "Employee"
	case PermissionLevelManager:
		return "Manager"
	case PermissionLevelAdmin:
		return "Administrator"
	}
	return "unknown"
}

type User struct {
	gorm.Model

	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	PermissionLevel PermissionLevel `json:"permission_level"`
	Email           string          `json:"email" gorm:"uniqueIndex"`
	Password        string          `json:"-"`

	OrganizationID uint         `json:"user_id"`
	Organization   Organization `json:"-"`
}

type Token struct {
	ID        string `json:"key" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint
	User      User
}

type UserWeekdayHours struct {
	gorm.Model

	UserID uint         `json:"user_id" gorm:"index:uwh_unique_oid_day,unique"`
	Day    time.Weekday `json:"day" gorm:"index:owh_unique_oid_day,unique"`

	Hours []UWHAvailability `json:"hours" gorm:"constraint:OnDelete:CASCADE"`
}
type UWHAvailability struct {
	gorm.Model

	UserWeekdayHoursID uint

	StartMinute uint `json:"start_minute"`
	EndMinute   uint `json:"end_minute"`
}

type UserOverrideHours struct {
	gorm.Model

	UserID uint      `json:"user_id" gorm:"index:uoh_unique_oid_date,unique"`
	Date   time.Time `json:"date" gorm:"index:uoh_unique_oid_date,unique"`

	Hours []UOHAvailability `json:"hours" gorm:"constraint:OnDelete:CASCADE"`
}
type UOHAvailability struct {
	gorm.Model

	UserOverrideHoursID uint

	StartMinute uint `json:"start_minute"`
	EndMinute   uint `json:"end_minute"`
}

func SetUserHours(db *gorm.DB, userid uint, hours OrganizationHours) {
	// step 1: delete all old hours
	var owh []UserWeekdayHours
	db.Unscoped().Clauses(clause.Returning{}).Where("user_id = ?", userid).Delete(&owh)
	for _, e := range owh {
		db.Unscoped().Clauses(clause.Returning{}).Where("user_weekday_hours_id = ?", e.ID).Delete(&UWHAvailability{})
	}
	var ooh []UserOverrideHours
	db.Unscoped().Clauses(clause.Returning{}).Where("user_id = ?", userid).Delete(&ooh)
	for _, e := range ooh {
		db.Unscoped().Clauses(clause.Returning{}).Where("user_weekday_hours_id = ?", e.ID).Delete(&UOHAvailability{})
	}

	// step 2: create new weekday hours
	monday := &UserWeekdayHours{
		UserID: userid,
		Day:    time.Monday,
	}
	db.Create(&monday)
	for _, e := range hours.Days.Monday {
		h := &UWHAvailability{
			UserWeekdayHoursID: monday.ID,
			StartMinute:        e.Start,
			EndMinute:          e.End,
		}
		db.Create(&h)
	}
	tuesday := &UserWeekdayHours{
		UserID: userid,
		Day:    time.Tuesday,
	}
	db.Create(&tuesday)
	for _, e := range hours.Days.Tuesday {
		h := &UWHAvailability{
			UserWeekdayHoursID: tuesday.ID,
			StartMinute:        e.Start,
			EndMinute:          e.End,
		}
		db.Create(&h)
	}
	wednesday := &UserWeekdayHours{
		UserID: userid,
		Day:    time.Wednesday,
	}
	db.Create(&wednesday)
	for _, e := range hours.Days.Wednesday {
		h := &UWHAvailability{
			UserWeekdayHoursID: wednesday.ID,
			StartMinute:        e.Start,
			EndMinute:          e.End,
		}
		db.Create(&h)
	}
	thursday := &UserWeekdayHours{
		UserID: userid,
		Day:    time.Thursday,
	}
	db.Create(&thursday)
	for _, e := range hours.Days.Thursday {
		h := &UWHAvailability{
			UserWeekdayHoursID: thursday.ID,
			StartMinute:        e.Start,
			EndMinute:          e.End,
		}
		db.Create(&h)
	}
	friday := &UserWeekdayHours{
		UserID: userid,
		Day:    time.Friday,
	}
	db.Create(&friday)
	for _, e := range hours.Days.Friday {
		h := &UWHAvailability{
			UserWeekdayHoursID: friday.ID,
			StartMinute:        e.Start,
			EndMinute:          e.End,
		}
		db.Create(&h)
	}
	saturday := &UserWeekdayHours{
		UserID: userid,
		Day:    time.Saturday,
	}
	db.Create(&saturday)
	for _, e := range hours.Days.Saturday {
		h := &UWHAvailability{
			UserWeekdayHoursID: saturday.ID,
			StartMinute:        e.Start,
			EndMinute:          e.End,
		}
		db.Create(&h)
	}
	sunday := &UserWeekdayHours{
		UserID: userid,
		Day:    time.Sunday,
	}
	db.Create(&sunday)
	for _, e := range hours.Days.Sunday {
		h := &UWHAvailability{
			UserWeekdayHoursID: sunday.ID,
			StartMinute:        e.Start,
			EndMinute:          e.End,
		}
		db.Create(&h)
	}

	for _, e := range ooh {
		override := &UserOverrideHours{
			UserID: userid,
			Date:   e.Date,
		}
		db.Create(&override)
		for _, f := range e.Hours {
			h := &UOHAvailability{
				UserOverrideHoursID: thursday.ID,
				StartMinute:         f.StartMinute,
				EndMinute:           f.EndMinute,
			}
			db.Create(&h)
		}
	}
}
func UWHToAvailability(owh UserWeekdayHours) []Availability {
	avail := make([]Availability, 0)
	for _, e := range owh.Hours {
		avail = append(avail, Availability{
			Start: e.StartMinute,
			End:   e.EndMinute,
		})
	}
	return avail
}
func UOHToAvailability(ooh UserOverrideHours) []Availability {
	avail := make([]Availability, 0)
	for _, e := range ooh.Hours {
		avail = append(avail, Availability{
			Start: e.StartMinute,
			End:   e.EndMinute,
		})
	}
	return avail
}
func GetUserHours(db *gorm.DB, userid uint) OrganizationHours {
	// get weekday hours
	var monday UserWeekdayHours
	db.Preload("Hours").Where("user_id = ? AND day = ?", userid, time.Monday).First(&monday)
	var tuesday UserWeekdayHours
	db.Preload("Hours").Where("user_id = ? AND day = ?", userid, time.Tuesday).First(&tuesday)
	var wednesday UserWeekdayHours
	db.Preload("Hours").Where("user_id = ? AND day = ?", userid, time.Wednesday).First(&wednesday)
	var thursday UserWeekdayHours
	db.Preload("Hours").Where("user_id = ? AND day = ?", userid, time.Thursday).First(&thursday)
	var friday UserWeekdayHours
	db.Preload("Hours").Where("user_id = ? AND day = ?", userid, time.Friday).First(&friday)
	var saturday UserWeekdayHours
	db.Preload("Hours").Where("user_id = ? AND day = ?", userid, time.Saturday).First(&saturday)
	var sunday UserWeekdayHours
	db.Preload("Hours").Where("user_id = ? AND day = ?", userid, time.Sunday).First(&sunday)

	// get override hours
	var overrides []UserOverrideHours
	db.Preload("Hours").Where("user_id = ?", userid).Find(&overrides)

	doh := make([]DayOverride, 0)
	for _, e := range overrides {
		doh = append(doh, DayOverride{
			Date:  e.Date,
			Hours: UOHToAvailability(e),
		})
	}

	hours := OrganizationHours{
		Days: WeekdayAvailability{
			Monday:    UWHToAvailability(monday),
			Tuesday:   UWHToAvailability(tuesday),
			Wednesday: UWHToAvailability(wednesday),
			Thursday:  UWHToAvailability(thursday),
			Friday:    UWHToAvailability(friday),
			Saturday:  UWHToAvailability(saturday),
			Sunday:    UWHToAvailability(sunday),
		},
		Overrides: doh,
	}

	return hours
}
