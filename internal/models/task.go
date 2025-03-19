package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Priority int

type Task struct {
	gorm.Model
	Title     string    `json:"title" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Priority  Priority  `json:"priority" gorm:"default:0" validate:"gte=0,lte=2"`
	Deadline  time.Time `json:"deadline" gorm:"not null"`
	Completed bool      `json:"completed" gorm:"default:false"`
	Id        int       `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	UserId    uuid.UUID `json:"user_id" gorm:"not null"`
}

const (
	Low Priority = iota
	Medium
	High
)

func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "Unknown"
	}
}

func (p *Priority) UnmarshalJSON(data []byte) error {
	var priorityStr string
	if err := json.Unmarshal(data, &priorityStr); err != nil {
		return err
	}
	switch priorityStr {
	case "Low":
		*p = Low
	case "Medium":
		*p = Medium
	case "High":
		*p = High
	default:
		return fmt.Errorf("invalid priority: %s", priorityStr)
	}
	return nil
}
