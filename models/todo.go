package models

import "time"

type Todo struct {
	ID   string `json:"id" bson:"_id"`
	Task string `json:"task" bson:"task"`
	Done bool   `json:"done" bson:"done"`
}

type DbTodo struct {
	ID          string    `json:"id" bson:"_id"`
	Task        string    `json:"task" bson:"task"`
	Done        bool      `json:"done" bson:"done"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	CompletedAt time.Time `json:"completedAt,omitempty" bson:"completedAt,omitempty"`
}
