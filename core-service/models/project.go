package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TenantID    bson.ObjectID `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	Name        string        `json:"name,omitempty" bson:"name,omitempty"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
