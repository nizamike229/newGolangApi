package models

import "time"

type TaskRequest struct {
	Title    string    `json:"title" validate:"required,min=1,max=100"`
	Priority Priority  `json:"priority" validate:"gte=0,lte=2"`
	Deadline time.Time `json:"deadline" validate:"required"`
}
