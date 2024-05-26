package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`

	Role string `json:"role"`
}

type Subscriber struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ListId uint           `json:"ListId"`
	List   SubscriberList `gorm:"foreignKey:ListId"`
}

type Admin struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name string `json:"name"`
}

type SubscriberList struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Subscribers []Subscriber `json:"subscribers" gorm:"foreignKey:ListId"`
}

// because the app can support multi account setup with the same UI, the user can switch between the accounts
type WhatsappBusinessAccount struct {
	gorm.Model
	UniqueId  uint      `json:"uniqueId" gorm:"primarykey,unique"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
