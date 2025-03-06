package middleware

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/khanirfan96/To-do-Fullstack-server/database"
	"github.com/khanirfan96/To-do-Fullstack-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTodo(c *fiber.Ctx) error {
	payload := getAllTasks(database.DB.TodoCollection)
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

func getAllTasks(coll *mongo.Collection) []primitive.M {
	cursor, err := coll.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var results []primitive.M
	for cursor.Next(context.Background()) {
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
	cursor.Close(context.Background())
	return results
}

func taskComplete(task string, newTask string, coll *mongo.Collection) error {
	id, err := primitive.ObjectIDFromHex(task)
	if err != nil {
		return fmt.Errorf("invalid ID: %v", err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true, "task": newTask}}

	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no document matched the given ID")
	}

	return nil
}

func insertOneTask(task models.ToDoList, coll *mongo.Collection) {
	insertResult, err := coll.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single record: ", insertResult.InsertedID)
}

func undoTask(task string, coll *mongo.Collection) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}
	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count", result.ModifiedCount)
}

func deleteOneTask(task string, coll *mongo.Collection) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	deletedID, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted Task: ", deletedID)
}

func deleteAllTask(coll *mongo.Collection) int64 {
	deletedAll, err := coll.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted All Tasks: ", deletedAll.DeletedCount)
	return deletedAll.DeletedCount
}
