package middleware

import (
	helper "github.com/khanirfan96/To-do-Fullstack-server/helpers"

	"github.com/gofiber/fiber/v2"
)

// Authz validates token and authorizes users
func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientToken := c.Get("token")
		if clientToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No Authorization header provided",
			})
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err,
			})
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		return c.Next()

	}
}
