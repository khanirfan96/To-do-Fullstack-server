package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	controller "github.com/khanirfan96/To-do-Fullstack-server/controller"
	"github.com/khanirfan96/To-do-Fullstack-server/database"
	"github.com/khanirfan96/To-do-Fullstack-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdatePassword(c *fiber.Ctx) error {
	userID := c.Params("id")

	var passwordUpdate models.UserPassword

	if err := c.BodyParser(&passwordUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = database.DB.UserCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	isValid, msg := controller.VerifyPassword(*user.Password, passwordUpdate.CurrentPassword)
	if !isValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": msg,
		})
	}

	// Hash the new password
	hashedPassword := controller.HashPassword(passwordUpdate.NewPassword)
	hashedPasswordPtr := &hashedPassword

	// Update the password in the database
	modifiedCount, err := updateUserPassword(userID, hashedPasswordPtr, database.DB.UserCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update password: %v", err),
		})
	}

	if modifiedCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Password update failed - no documents modified",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      userID,
		"message": "Password updated successfully!",
		"updated": modifiedCount,
		"status":  fiber.StatusOK,
	})
}

func updateUserPassword(id string, hashedPassword *string, userColl *mongo.Collection) (int64, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"password": hashedPassword,
		},
	}

	result, err := userColl.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}
