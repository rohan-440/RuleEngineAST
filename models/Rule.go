package models

import "time"

type Rule struct {
	Id        uint      `json:"id" gorm:"primary_key"`
	Rule      string    `json:"rule"`
	CreatedAt time.Time `json:"createdAt"`
}
