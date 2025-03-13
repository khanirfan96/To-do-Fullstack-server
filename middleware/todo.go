package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khanirfan96/To-do-Fullstack-server/database"
	"github.com/khanirfan96/To-do-Fullstack-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTodo(c *fiber.Ctx) error {
	uid := c.Locals("Uid").(string)
	payload := getAllTasks(database.DB.TodoCollection, uid)
	return c.JSON(payload)
}

func CreateTodo(c *fiber.Ctx) error {
	var task models.ToDoList
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	insertOneTask(task, database.DB.TodoCollection)
	return c.JSON(task)
}

func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		NewTask string `json:"task"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := taskComplete(id, body.NewTask, database.DB.TodoCollection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update task: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"id":      id,
		"message": "Task updated successfully",
	})
}

func UndoTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	undoTask(id, database.DB.TodoCollection)
	return c.JSON(id)
}

func DeleteOneTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	deleteOneTask(id, database.DB.TodoCollection)
	return c.JSON(id)
}

func DeleteAllTodo(c *fiber.Ctx) error {
	count := deleteAllTask(database.DB.TodoCollection)
	return c.JSON(count)
}

func getAllTasks(coll *mongo.Collection, id string) []primitive.M {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	filter := bson.M{"user_id": id}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []primitive.M
	for cursor.Next(ctx) {
		var result bson.M
		e := cursor.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	defer cancel()
	return results
}

func taskComplete(task string, newTask string, coll *mongo.Collection) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	id, err := primitive.ObjectIDFromHex(task)
	defer cancel()
	if err != nil {
		return fmt.Errorf("invalid ID: %v", err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true, "task": newTask}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no document matched the given ID")
	}

	return nil
}

func insertOneTask(task models.ToDoList, coll *mongo.Collection) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	insertResult, err := coll.InsertOne(ctx, task)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	fmt.Println("Inserted a single record: ", insertResult.InsertedID)
}

func undoTask(task string, coll *mongo.Collection) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	fmt.Println("Modified count", result.ModifiedCount)
}

func deleteOneTask(task string, coll *mongo.Collection) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	deletedID, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	fmt.Println("Deleted Task: ", deletedID)
}

func deleteAllTask(coll *mongo.Collection) int64 {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	deletedAll, err := coll.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	fmt.Println("Deleted All Tasks: ", deletedAll.DeletedCount)
	return deletedAll.DeletedCount
}
