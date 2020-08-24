package models

import "github.com/jinzhu/gorm"

// Gallery is the underlying data structure that serves as the model for users to view galleries
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}
