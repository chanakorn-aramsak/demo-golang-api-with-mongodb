package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book represents a book model
type Book struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title  string             `json:"title" bson:"title,omitempty"`
	Author string             `json:"author" bson:"author,omitempty"`
}
