package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Product struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	ImageLink   string        `json:"image_link" bson:"image_link"`
	Price       float64       `json:"price" bson:"price"`
	Discount    float64       `json:"discount" bson:"discount"`
}
