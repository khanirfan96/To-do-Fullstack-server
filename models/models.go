package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=8"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}

type UserPassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type Gym struct {
	ID        primitive.ObjectID `bson:"_id"`
	Monday    *string            `json:"monday"`
	Tuesday   *string            `json:"tuesday"`
	Wednesday *string            `json:"wednesday"`
	Thursday  *string            `json:"thursday"`
	Friday    *string            `json:"friday"`
	Saturday  *string            `json:"saturday"`
	Sunday    *string            `json:"sunday"`
}
