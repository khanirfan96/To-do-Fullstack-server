package controllers

import (
	"context"
	"fmt"
	"log"

	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/khanirfan96/To-do-Fullstack-server/database"

	helper "github.com/khanirfan96/To-do-Fullstack-server/helpers"
	"github.com/khanirfan96/To-do-Fullstack-server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "login or passowrd is incorrect"
		check = false
	}

	return check, msg
}

// CreateUser is the api used to tget a single user
func SignUp() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})

		}

		count, err := database.DB.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while checking for the email"})
		}

		if count > 0 {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "This email already exists"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = database.DB.UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error occured while checking for the phone number"})
		}

		if count > 0 {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "This phone number already exists"})
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := database.DB.UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "User item was not created"})

		}

		return c.Status(http.StatusOK).JSON(resultInsertionNumber)

	}
}

// controller/userController.go
func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Check if UserCollection is initialized
		if database.DB.UserCollection == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not properly initialized"})
		}

		var user models.User
		var foundUser models.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if user.Email == nil || user.Password == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
		}

		err := database.DB.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
		}

		token, refreshToken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate tokens"})
		}

		if err := helper.UpdateAllTokens(token, refreshToken, foundUser.User_id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to update tokens: %v", err),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": foundUser})
	}
}
