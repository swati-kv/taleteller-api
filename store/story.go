package store

import (
	"context"
	"time"
)

type Story struct {
	StoryID     string    `db:"id" json:"story_id,omitempty"`
	Name        string    `db:"name" json:"name,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
	Mood        string    `db:"mood" json:"mood,omitempty"`
	Category    string    `db:"category" json:"category,omitempty"`
	CustomerID  string    `db:"customer_id" json:"customer_id,omitempty"`
	Status      string    `db:"status" json:"status,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type StoryStorer interface {
	Create(ctx context.Context, createRequest Story) (err error)
	List(ctx context.Context) (stories []Story, err error)
}
