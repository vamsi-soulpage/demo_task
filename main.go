package main

import (
	"demo/crud"
	"demo/db"
	"errors"
	_ "errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg0Nzc4NzMsInVzZXJuYW1lIjoidmFtc2kifQ.Vszoo66Wl1B9BUILuFBp2ydFiBKvDneYK580Gzy_kos
var jwtSecret = []byte("soulpage")
var User db.User

func Register(c *gin.Context) {
	var Sign db.User
	Sign.Username = c.Request.PostFormValue("username")
	Sign.Email = c.Request.PostFormValue("email")
	Sign.Password = c.Request.PostFormValue("password")

	// Create the user record
	user := db.User{
		Username: Sign.Username,
		Email:    Sign.Email,
		Password: Sign.Password,
	}

	var existingUser db.User

	if err := db.DBASE.Where("username = ? OR email = ?", Sign.Username, Sign.Email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error"})
			return
		}
	} else {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	if err := db.DBASE.Create(&user).Error; err != nil {
		fmt.Println("error creating user in the database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var A db.User
	A.Username = c.Request.PostFormValue("username")
	A.Password = c.Request.PostFormValue("password")

	var existingUser db.User
	fmt.Println(A.Username)
	if err := db.DBASE.Where("username = ?", A.Username).First(&existingUser).Error; err != nil {
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

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "message": "Login successful"})
}

func main() {
	db.Connect()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	gin.Recovery()
	gin.Logger()
	router.POST("/signup", Register)
	router.POST("/login", Login)
	router.POST("/create", crud.CreateUserData)
	router.POST("/getuser", crud.GetUserData)
	router.POST("/update", crud.UpdateUserDate)
	router.POST("/delete", crud.DeleteUserData)
	err := router.Run(":9001")
	if err != nil {
		fmt.Println("error at port", err.Error())
	}
}
