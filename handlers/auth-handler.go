package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dp487/legendary-succotash/database"
	"github.com/dp487/legendary-succotash/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = os.Getenv("TOKEN_SECRET")

// Centralized function for JSON responses
func respond(c *fiber.Ctx, status int, message string, data interface{}) error {
	response := models.BuildResponse(http.StatusText(status), message, data, "")
	return c.Status(status).JSON(response)
}

// Register handler
func HandleRegister(c *fiber.Ctx, db *database.Database) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return respond(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{Username: data["username"]}
	user.SetPassword(password)

	var existingUser models.User
	result := db.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		return respond(c, fiber.StatusBadRequest, "Username already exists", nil)
	}

	if u := db.DB.Create(&user); u.Error != nil {
		return respond(c, fiber.StatusInternalServerError, "Failed to create new user", u.Error.Error())
	}

	return respond(c, fiber.StatusCreated, "User created successfully", user)
}

// Login handler
func HandleLogin(c *fiber.Ctx, db *database.Database) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return respond(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	var user models.User
	dbc := db.DB.Where("username = ?", data["username"]).First(&user)
	if dbc.RowsAffected == 0 {
		return respond(c, fiber.StatusNotFound, "User not found", nil)
	}

	if err := bcrypt.CompareHashAndPassword(user.GetPassword(), []byte(data["password"])); err != nil {
		return respond(c, fiber.StatusBadRequest, "Incorrect password", nil)
	}

	var session models.UserSessions
	qr := db.DB.Where("username = ?", user.Username).First(&session)
	if qr.Error == nil {
		return respond(c, fiber.StatusBadRequest, "User already logged in", nil)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    user.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 1 Day
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, "Could not login", err.Error())
	}

	session = models.UserSessions{
		Username: user.Username,
		Token:    token,
	}

	if result := db.DB.Create(&session); result.Error != nil {
		return respond(c, fiber.StatusInternalServerError, "Error creating session", result.Error.Error())
	}

	c.Set("Authorization", "Bearer "+token)
	return respond(c, fiber.StatusOK, "Logged in", user)
}

// Logout handler
func HandleLogout(c *fiber.Ctx, db *database.Database) error {
	authorized, err := isAuthenticated(c, db)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
	}

	if !authorized {
		return respond(c, fiber.StatusUnauthorized, "Not authorized", nil)
	}

	requestToken := c.Get("Authorization")
	requestToken = strings.TrimPrefix(requestToken, "Bearer ")

	_, claims, err := getJwtTokenAndClaims(requestToken)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, "Error logging out", err.Error())
	}

	if qr := db.DB.Where("username = ?", claims.Issuer).Delete(&models.UserSessions{}); qr.Error != nil {
		return respond(c, fiber.StatusInternalServerError, "Error logging out", qr.Error.Error())
	}

	return respond(c, fiber.StatusOK, "Logged out", nil)
}

// Check authentication handler
func HandleIsAuthenticated(c *fiber.Ctx, db *database.Database) error {
	authorized, err := isAuthenticated(c, db)
	if err != nil {
		return respond(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
	}

	if !authorized {
		return respond(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	return respond(c, fiber.StatusOK, "Authorized", nil)
}

// Authentication check function
func isAuthenticated(c *fiber.Ctx, db *database.Database) (bool, error) {
	requestToken := c.Get("Authorization")
	requestToken = strings.TrimPrefix(requestToken, "Bearer ")

	jwtToken, claims, err := getJwtTokenAndClaims(requestToken)
	if err != nil {
		return false, err
	}

	var userSession models.UserSessions
	qr := db.DB.Where("username = ?", claims.Issuer).First(&userSession)

	if qr.Error != nil {
		return false, qr.Error
	}

	if qr.RowsAffected == 0 {
		return false, nil
	}

	if userSession.Token != jwtToken.Raw {
		return false, nil
	}

	storedTime := claims.ExpiresAt.Time
	currentTime := time.Now()

	if currentTime.After(storedTime) {
		return false, nil
	}

	return true, nil
}

// JWT token and claims extraction function
func getJwtTokenAndClaims(token string) (*jwt.Token, *jwt.RegisteredClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims := jwtToken.Claims.(*jwt.RegisteredClaims)
	return jwtToken, claims, nil
}
