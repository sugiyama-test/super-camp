package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Checklist struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChecklistItem struct {
	ID          int64     `json:"id"`
	ChecklistID int64    `json:"checklist_id"`
	Name        string   `json:"name"`
	IsChecked   bool     `json:"is_checked"`
	Quantity    int      `json:"quantity"`
	SortOrder   int      `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Layout struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FireLog struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Date            time.Time `json:"date"`
	Location        string    `json:"location"`
	WoodType        string    `json:"wood_type"`
	DurationMinutes int       `json:"duration_minutes"`
	Notes           string    `json:"notes"`
	Temperature     *float64  `json:"temperature,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MealPlan struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	MealType  string    `json:"meal_type"`
	Servings  int       `json:"servings"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Gear struct {
	ID          int64    `json:"id"`
	UserID      int64    `json:"user_id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Brand       string   `json:"brand"`
	WeightGrams *float64 `json:"weight_grams,omitempty"`
	Notes       string   `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Campsite struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Notes     string   `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
