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

func GetRecipe(c *fiber.Ctx) error {
	payload := getAllCalories(database.DB.CalorieCollection)
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
	insertOneRecipe(recipe, database.DB.CalorieCollection)
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
	count := deleteAllRecipe(database.DB.CalorieCollection)
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
	deleteOneRecipe(id, database.DB.CalorieCollection)
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

	var request models.CalorieTracker

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	modifiedCount, err := updateRecipe(id, request, database.DB.CalorieCollection)

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

	modifiedIngredient, err := updateIngredients(id, ingredients, database.DB.CalorieCollection)

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
