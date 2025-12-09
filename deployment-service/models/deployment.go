package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DeploymentStatus string

const (
	StatusPending DeploymentStatus = "Pending"
	StatusRunning DeploymentStatus = "Running"
	StatusFailed  DeploymentStatus = "Failed"
)

type Deployment struct {
	ID        bson.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	ServiceID bson.ObjectID    `json:"service_id,omitempty" bson:"service_id,omitempty"`
	Status    DeploymentStatus `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt time.Time        `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time        `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
