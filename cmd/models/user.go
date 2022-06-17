package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id"`
	CreatedBy  string `json:"-"`
	UserName   string `gorm:"size:255;not null" json:"userName"`
	Passwords  string `gorm:"size:255;not null" json:"-"`
	UserType   string `gorm:"size:1;not null" json:"userType"`
	UserStatus string `gorm:"size:1;not null" json:"userStatus"`
	Code       string `gorm:"size:100;not null" json:"-"`
	FirstName  string `gorm:"size:100;" json:"firstName"`
	LastName   string `gorm:"size:100;" json:"lastName"`
	IdCard     string `gorm:"size:13;" json:"-"`
	Province   string `gorm:"size:100;" json:"province"`
	Amphure    string `gorm:"size:100;" json:"-"`
	District   string `gorm:"size:100;" json:"-"`
	PdpaCheck  string `gorm:"size:1;" json:"pdpaCheck"`
}
