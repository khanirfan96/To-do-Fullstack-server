package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToDoList struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Task   string             `json:"task,omitempty"`
	Status bool               `json:"status,omitempty"`
}

type CalorieTracker struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Dish        *string            `json:"dish"`
	Ingredients *string            `json:"ingredients"`
	Calories    *int64             `json:"calories"`
	Fat         *int64             `json:"fat"`
}
