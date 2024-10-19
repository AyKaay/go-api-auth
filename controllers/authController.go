package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"user-auth-api/initializers"
	"user-auth-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var authInput models.AuthInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := initializers.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not load users",
			"details": err.Error(),
		})
		return
	}

	for _, user := range users {
		if user.Username == authInput.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
			return
		}
	}

	// hash input password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := models.User{
		Username: authInput.Username,
		Password: string(passwordHash),
	}

	if err := initializers.AddUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newUser})
}

func Login(c *gin.Context) {

	var authInput models.AuthInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := initializers.LoadUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not load users",
			"details": err.Error(),
		})
		return
	}

	var foundUser *models.User
	// Check if the username exists and find the user
	for _, user := range users {
		if user.Username == authInput.Username {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	// generate jwt token
	expire, err := strconv.Atoi(os.Getenv("EXPIRE"))
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  foundUser.ID,
		"exp": time.Now().Add(time.Hour * time.Duration(expire)).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

var tokenBlacklist = make(map[string]bool)

func Logout(c *gin.Context) {
	// get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	tokenString := authToken[1]

	// Invalidate the token by adding it to the blacklist
	tokenBlacklist[tokenString] = true
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func CreateTransaction(c *gin.Context) {
	var paymentInput models.PaymentInput
	if err := c.ShouldBindJSON(&paymentInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	transactions, err := initializers.LoadTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load transactions"})
		return
	}

	newTransaction := models.Transaction{
		ID:     uuid.New(),
		Amount: paymentInput.Amount,
		Type:   paymentInput.Type,
		UserID: userID.(uuid.UUID),
	}

	transactions = append(transactions, newTransaction)

	if err := initializers.WriteTransactions(transactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newTransaction})
}
