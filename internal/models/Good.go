package models

import (
	"fmt"
	"time"
)

type Good struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (g *Good) Key() string {
	return fmt.Sprintf("id=%d&projectId=%d", g.ID, g.ProjectID)
}

type UpdateGoodBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (g Good) GetLogEvent() GoodLogEvent {
	return GoodLogEvent{
		Good: g,

		EventTime: time.Now(),
	}
}

type GoodLogEvent struct {
	Good

	EventTime time.Time
}
