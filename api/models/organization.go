package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Organization struct {
	gorm.Model

	Name        string `json:"name" gorm:"unique"`
	EmailDomain string `json:"email_domain" gorm:"unique"`

	WeekdayHours  []OrganizationWeekdayHours  `json:"hours" gorm:"constraint:OnDelete:CASCADE"`
	OverrideHours []OrganizationOverrideHours `json:"override_hours" gorm:"constraint:OnDelete:CASCADE"`
}

type OrganizationWeekdayHours struct {
	gorm.Model

	OrganizationID uint         `json:"organization_id" gorm:"index:owh_unique_oid_day,unique"`
	Day            time.Weekday `json:"day" gorm:"index:owh_unique_oid_day,unique"`

	Hours []OWHAvailability `json:"hours" gorm:"constraint:OnDelete:CASCADE"`
}
type OWHAvailability struct {
	gorm.Model

	OrganizationWeekdayHoursID uint

	StartMinute uint `json:"start_minute"`
	EndMinute   uint `json:"end_minute"`
}

type OrganizationOverrideHours struct {
	gorm.Model

	OrganizationID uint      `json:"organization_id" gorm:"index:ooh_unique_oid_date,unique"`
	Date           time.Time `json:"date" gorm:"index:ooh_unique_oid_date,unique"`

	Hours []OOHAvailability `json:"hours" gorm:"constraint:OnDelete:CASCADE"`
}
type OOHAvailability struct {
	gorm.Model

	OrganizationOverrideHoursID uint

	StartMinute uint `json:"start_minute"`
	EndMinute   uint `json:"end_minute"`
}

type Availability struct {
	Start uint `json:"start"`
	End   uint `json:"end"`
}
type DayOverride struct {
	Date  time.Time      `json:"date"`
	Hours []Availability `json:"hours"`
}
type WeekdayAvailability struct {
	Sunday    []Availability `json:"sunday"`
	Monday    []Availability `json:"monday"`
	Tuesday   []Availability `json:"tuesday"`
	Wednesday []Availability `json:"wednesday"`
	Thursday  []Availability `json:"thursday"`
	Friday    []Availability `json:"friday"`
	Saturday  []Availability `json:"saturday"`
}
type OrganizationHours struct {
	Days      WeekdayAvailability `json:"days"`
	Overrides []DayOverride       `json:"overrides"`
}

func HMToMinute(hour uint, minute uint) uint {
	return hour*60 + minute
}

func SetOrgHours(db *gorm.DB, orgid uint, hours OrganizationHours) {
	// step 1: delete all old hours
	var owh []OrganizationWeekdayHours
	db.Unscoped().Clauses(clause.Returning{}).Where("organization_id = ?", orgid).Delete(&owh)
	for _, e := range owh {
		db.Unscoped().Clauses(clause.Returning{}).Where("organization_weekday_hours_id = ?", e.ID).Delete(&OWHAvailability{})
	}
	var ooh []OrganizationOverrideHours
	db.Unscoped().Clauses(clause.Returning{}).Where("organization_id = ?", orgid).Delete(&ooh)
	for _, e := range ooh {
		db.Unscoped().Clauses(clause.Returning{}).Where("organization_weekday_hours_id = ?", e.ID).Delete(&OOHAvailability{})
	}

	// step 2: create new weekday hours
	monday := &OrganizationWeekdayHours{
		OrganizationID: orgid,
		Day:            time.Monday,
	}
	db.Create(&monday)
	for _, e := range hours.Days.Monday {
		h := &OWHAvailability{
			OrganizationWeekdayHoursID: monday.ID,
			StartMinute:                e.Start,
			EndMinute:                  e.End,
		}
		db.Create(&h)
	}
	tuesday := &OrganizationWeekdayHours{
		OrganizationID: orgid,
		Day:            time.Tuesday,
	}
	db.Create(&tuesday)
	for _, e := range hours.Days.Tuesday {
		h := &OWHAvailability{
			OrganizationWeekdayHoursID: tuesday.ID,
			StartMinute:                e.Start,
			EndMinute:                  e.End,
		}
		db.Create(&h)
	}
	wednesday := &OrganizationWeekdayHours{
		OrganizationID: orgid,
		Day:            time.Wednesday,
	}
	db.Create(&wednesday)
	for _, e := range hours.Days.Wednesday {
		h := &OWHAvailability{
			OrganizationWeekdayHoursID: wednesday.ID,
			StartMinute:                e.Start,
			EndMinute:                  e.End,
		}
		db.Create(&h)
	}
	thursday := &OrganizationWeekdayHours{
		OrganizationID: orgid,
		Day:            time.Thursday,
	}
	db.Create(&thursday)
	for _, e := range hours.Days.Thursday {
		h := &OWHAvailability{
			OrganizationWeekdayHoursID: thursday.ID,
			StartMinute:                e.Start,
			EndMinute:                  e.End,
		}
		db.Create(&h)
	}
	friday := &OrganizationWeekdayHours{
		OrganizationID: orgid,
		Day:            time.Friday,
	}
	db.Create(&friday)
	for _, e := range hours.Days.Friday {
		h := &OWHAvailability{
			OrganizationWeekdayHoursID: friday.ID,
			StartMinute:                e.Start,
			EndMinute:                  e.End,
		}
		db.Create(&h)
	}
	saturday := &OrganizationWeekdayHours{
		OrganizationID: orgid,
		Day:            time.Saturday,
	}
	db.Create(&saturday)
	for _, e := range hours.Days.Saturday {
		h := &OWHAvailability{
			OrganizationWeekdayHoursID: saturday.ID,
			StartMinute:                e.Start,
			EndMinute:                  e.End,
		}
		db.Create(&h)
	}
	sunday := &OrganizationWeekdayHours{
		OrganizationID: orgid,
		Day:            time.Sunday,
	}
	db.Create(&sunday)
	for _, e := range hours.Days.Sunday {
		h := &OWHAvailability{
			OrganizationWeekdayHoursID: sunday.ID,
			StartMinute:                e.Start,
			EndMinute:                  e.End,
		}
		db.Create(&h)
	}

	for _, e := range ooh {
		override := &OrganizationOverrideHours{
			OrganizationID: orgid,
			Date:           e.Date,
		}
		db.Create(&override)
		for _, f := range e.Hours {
			h := &OOHAvailability{
				OrganizationOverrideHoursID: thursday.ID,
				StartMinute:                 f.StartMinute,
				EndMinute:                   f.EndMinute,
			}
			db.Create(&h)
		}
	}
}
func OWHToAvailability(owh OrganizationWeekdayHours) []Availability {
	avail := make([]Availability, 0)
	for _, e := range owh.Hours {
		avail = append(avail, Availability{
			Start: e.StartMinute,
			End:   e.EndMinute,
		})
	}
	return avail
}
func OOHToAvailability(ooh OrganizationOverrideHours) []Availability {
	avail := make([]Availability, 0)
	for _, e := range ooh.Hours {
		avail = append(avail, Availability{
			Start: e.StartMinute,
			End:   e.EndMinute,
		})
	}
	return avail
}
func GetOrgHours(db *gorm.DB, orgid uint) OrganizationHours {
	// get weekday hours
	var monday OrganizationWeekdayHours
	db.Preload("Hours").Where("organization_id = ? AND day = ?", orgid, time.Monday).First(&monday)
	var tuesday OrganizationWeekdayHours
	db.Preload("Hours").Where("organization_id = ? AND day = ?", orgid, time.Tuesday).First(&tuesday)
	var wednesday OrganizationWeekdayHours
	db.Preload("Hours").Where("organization_id = ? AND day = ?", orgid, time.Wednesday).First(&wednesday)
	var thursday OrganizationWeekdayHours
	db.Preload("Hours").Where("organization_id = ? AND day = ?", orgid, time.Thursday).First(&thursday)
	var friday OrganizationWeekdayHours
	db.Preload("Hours").Where("organization_id = ? AND day = ?", orgid, time.Friday).First(&friday)
	var saturday OrganizationWeekdayHours
	db.Preload("Hours").Where("organization_id = ? AND day = ?", orgid, time.Saturday).First(&saturday)
	var sunday OrganizationWeekdayHours
	db.Preload("Hours").Where("organization_id = ? AND day = ?", orgid, time.Sunday).First(&sunday)

	// get override hours
	var overrides []OrganizationOverrideHours
	db.Preload("Hours").Where("organization_id = ?", orgid).Find(&overrides)

	doh := make([]DayOverride, 0)
	for _, e := range overrides {
		doh = append(doh, DayOverride{
			Date:  e.Date,
			Hours: OOHToAvailability(e),
		})
	}

	hours := OrganizationHours{
		Days: WeekdayAvailability{
			Monday:    OWHToAvailability(monday),
			Tuesday:   OWHToAvailability(tuesday),
			Wednesday: OWHToAvailability(wednesday),
			Thursday:  OWHToAvailability(thursday),
			Friday:    OWHToAvailability(friday),
			Saturday:  OWHToAvailability(saturday),
			Sunday:    OWHToAvailability(sunday),
		},
		Overrides: doh,
	}

	return hours
}
