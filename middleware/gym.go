package middleware

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khanirfan96/To-do-Fullstack-server/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetGym(c *fiber.Ctx) error {
	payload := getAllGymSchedule(database.DB.GymCollection)
	return c.Status(fiber.StatusOK).JSON(payload)
}

func getAllGymSchedule(gymcoll *mongo.Collection) []primitive.M {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	cursor, err := gymcoll.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var gym []primitive.M
	for cursor.Next(ctx) {
		var gymschedule bson.M
		e := cursor.Decode(&gymschedule)
		if e != nil {
			log.Fatal(e)
		}
		gym = append(gym, gymschedule)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	defer cancel()
	return gym
}
