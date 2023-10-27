package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var DBASE *gorm.DB

var jwtSecret = []byte("soulpage")

type User struct {
	gorm.Model
	ID       uint   `gorm:"primary_key"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Connect() *gorm.DB {
	dbURL := "host=localhost user=postgres password=vamsi123 dbname=soulpage port=5432 sslmode=disable"
	Database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	} else {
		fmt.Println("connected to the database successfully")
	}

	err = Database.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err.Error())
	}

	DBASE = Database
	return DBASE
}

func Register(c *gin.Context) {
	var Sign User
	Sign.Username = c.Request.PostFormValue("username")
	Sign.Email = c.Request.PostFormValue("email")
	Sign.Password = c.Request.PostFormValue("password")

	// Create the user record
	user := User{
		Username: Sign.Username,
		Email:    Sign.Email,
		Password: Sign.Password,
	}

	var existingUser User

	if err := DBASE.Where("username = ? OR email = ?", Sign.Username, Sign.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	if err := DBASE.Create(&user).Error; err != nil {
		fmt.Println("error creating user in the database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var A User
	A.Username = c.Request.PostFormValue("username")
	A.Password = c.Request.PostFormValue("password")

	var existingUser User

	if err := DBASE.Where("username = ?", A.Username).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user with that name"})
		return
	}

	if existingUser.Password != A.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate and send a JWT token upon successful login
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = existingUser.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	//c.JSON(http.StatusOK, gin.H{"token": tokenString, "message": "Login successful"})
	extractedUsername, err := ExtractUsernameFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "message": "Login successful", "username": extractedUsername})
}

func ExtractUsernameFromToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Use the same jwtSecret that was used to sign the token
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Invalid token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("Username not found in token claims")
	}

	return username, nil
}

func main() {
	DBASE = Connect()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	gin.Recovery()
	gin.Logger()
	router.POST("/signup", Register)
	router.POST("/login", Login)
	err := router.Run(":9001")
	if err != nil {
		fmt.Println("error at port", err.Error())
	}
}
