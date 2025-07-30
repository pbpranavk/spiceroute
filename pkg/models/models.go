package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Preferences []Preference `gorm:"foreignKey:UserID" json:"preferences,omitempty"`
	Feedback    []Feedback   `gorm:"foreignKey:UserID" json:"feedback,omitempty"`
}

// Preference represents user dietary preferences
type Preference struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Preferences
	Cuisines   []string `gorm:"type:text[]" json:"cuisines"`
	Allergies  []string `gorm:"type:text[]" json:"allergies"`
	BudgetWeek float64  `json:"budget_week"`
	Spicy      bool     `json:"spicy"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Recipe represents a recipe in the system
type Recipe struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name          string         `gorm:"not null" json:"name"`
	Cuisine       string         `json:"cuisine"`
	PrepMinutes   int32          `json:"prep_minutes"`
	Calories      int32          `json:"calories"`
	Ingredients   []string       `gorm:"type:text[]" json:"ingredients"`
	Cost          float64        `json:"cost"`
	ShelfLifeDays int32          `json:"shelf_life_days"`
	Tags          []string       `gorm:"type:text[]" json:"tags"`
	Nutrition     string         `json:"nutrition"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Feedback []Feedback `gorm:"foreignKey:DishID" json:"feedback,omitempty"`
}

// Feedback represents user feedback on dishes
type Feedback struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          string         `gorm:"type:uuid;not null" json:"user_id"`
	DishID          string         `gorm:"type:uuid;not null" json:"dish_id"`
	Rating          int32          `json:"rating"`
	Skipped         bool           `json:"skipped"`
	SubstitutedWith string         `json:"substituted_with"`
	Comment         string         `json:"comment"`
	CookedAt        time.Time      `json:"cooked_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Recipe Recipe `gorm:"foreignKey:DishID" json:"recipe,omitempty"`
}

// TableName specifies the table name for Feedback
func (Feedback) TableName() string {
	return "feedback"
}

// TableName specifies the table name for Preference
func (Preference) TableName() string {
	return "preferences"
}
