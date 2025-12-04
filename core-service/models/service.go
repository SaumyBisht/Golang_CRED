package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProjectID   bson.ObjectID `json:"project_id,omitempty" bson:"project_id,omitempty"`
	Name        string        `json:"name,omitempty" bson:"name,omitempty"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
