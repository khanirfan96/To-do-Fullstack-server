package middleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/khanirfan96/To-do-Fullstack.git/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var collection *mongo.Collection

var (
	todoCollection    *mongo.Collection
	calorieCollection *mongo.Collection
)

func init() {
	loadTheEnv()
	createDBINstance()
}

func loadTheEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in dotEnv File")
	}
}

func createDBINstance() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	// collectionName := os.Getenv("DB_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	todoCollection = client.Database(dbName).Collection("todolist")
	calorieCollection = client.Database(dbName).Collection("calorietracker")

	fmt.Println("Connected to mongoDB!")
	fmt.Println("Collection instances created...")
	// collection = client.Database(dbName).Collection(collectionName)
	fmt.Printf("Todo Collection: %v\n", todoCollection.Name())
	fmt.Printf("Calorie Collection: %v\n", calorieCollection.Name())
	// fmt.Println("Collection instance created...", collection)
}

func GetTodo(c *fiber.Ctx) error {
	payload := getAllTasks(todoCollection)
	return c.JSON(payload)
}

func CreateTodo(c *fiber.Ctx) error {
	var task models.ToDoList
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	insertOneTask(task, todoCollection)
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

	if err := taskComplete(id, body.NewTask, todoCollection); err != nil {
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
	undoTask(id, todoCollection)
	return c.JSON(id)
}

func DeleteOneTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	deleteOneTask(id, todoCollection)
	return c.JSON(id)
}

func DeleteAllTodo(c *fiber.Ctx) error {
	count := deleteAllTask(todoCollection)
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
		fmt.Println("result", result)
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

//**************************CalorieTracker Api Methods**********************************************

func GetRecipe(c *fiber.Ctx) error {
	payload := getAllCalories(calorieCollection)
	return c.Status(fiber.StatusOK).JSON(payload)
}

func getAllCalories(calCol *mongo.Collection) []primitive.M {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	cursor, err := calCol.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var calories []primitive.M
	for cursor.Next(ctx) {
		var calorie bson.M
		e := cursor.Decode(&calorie)
		if e != nil {
			log.Fatal(e)
		}
		calories = append(calories, calorie)

	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	defer cancel()
	return calories
}

func CreateRecipe(c *fiber.Ctx) error {
	var recipe models.CalorieTracker
	if err := c.BodyParser(&recipe); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse json",
		})
	}
	insertOneRecipe(recipe, calorieCollection)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Recipe created successfully",
		"id":      recipe,
	})
}

func insertOneRecipe(recipe models.CalorieTracker, recipeColl *mongo.Collection) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	insertCalorieResult, err := recipeColl.InsertOne(ctx, recipe)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	fmt.Println("Inserted a Recipe ", insertCalorieResult.InsertedID)
}

func DeleteAllRecipe(c *fiber.Ctx) error {
	count := deleteAllRecipe(calorieCollection)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "All Entries Deleted",
		"Count":   count,
	})
}

func deleteAllRecipe(recipeColl *mongo.Collection) int64 {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	deletedAll, err := recipeColl.DeleteMany(ctx, bson.D{{}})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Deleted All Recipes ", deletedAll.DeletedCount)
	defer cancel()
	return deletedAll.DeletedCount
}

func DeleteOneRecipe(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println("id", id)
	deleteOneRecipe(id, calorieCollection)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": fmt.Sprintf("Deleted entry with ID: %s", id),
		"ID":      id,
	})
}

func deleteOneRecipe(id string, recipeColl *mongo.Collection) {
	ids, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": ids}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	deletedId, err := recipeColl.DeleteOne(ctx, filter)

	defer cancel()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Deleted Recipe ", deletedId)
}

func UpdateRecipe(c *fiber.Ctx) error {
	id := c.Params("id")

	var body models.CalorieTracker

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	modifiedCount, err := updateRecipe(id, body, calorieCollection)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update recipe: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      id,
		"message": "Recipe updated successfully",
		"updated": modifiedCount,
	})
}

func updateRecipe(id string, body models.CalorieTracker, recipeColl *mongo.Collection) (int64, error) {
	recipeId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("invalid recipe ID format")
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"dish":        body.Dish,
			"ingredients": body.Ingredients,
			"calories":    body.Calories,
			"fat":         body.Fat,
		},
	}
	result, err := recipeColl.UpdateOne(ctx, bson.M{"_id": recipeId}, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func UpdateIngredeints(c *fiber.Ctx) error {
	id := c.Params("id")

	var ingredients models.CalorieTracker

	if err := c.BodyParser(&ingredients); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	modifiedIngredient, err := updateIngredients(id, ingredients, calorieCollection)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update recipe: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      id,
		"message": "Ingredients updated successfully",
		"updated": modifiedIngredient,
	})

}

func updateIngredients(id string, ingredients models.CalorieTracker, recipeColl *mongo.Collection) (int64, error) {
	ingredientId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("invalid ingredient id%v", ingredientId)

	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"ingredients": ingredients.Ingredients,
		},
	}

	result, err := recipeColl.UpdateOne(ctx, bson.M{"_id": ingredientId}, update)

	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}
