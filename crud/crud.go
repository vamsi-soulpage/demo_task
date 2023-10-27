package crud

import (
	"demo/auth"
	"demo/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUserData(c *gin.Context) {
	//only admin user can create the newuser
	adminuser, err := auth.ExtractUsernameFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	fmt.Println("admin user trying to create a new user...>", adminuser)

	var userData db.User
	userData.Username = c.Request.PostFormValue("username")
	userData.Email = c.Request.PostFormValue("email")
	userData.Password = c.Request.PostFormValue("password")

	// Create the user data in the database
	if err := db.DBASE.Create(&userData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User data created"})
}

func GetUserData(c *gin.Context) {
	// admin user has access to see the data
	Adminuser, err := auth.ExtractUsernameFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	fmt.Println("admin user trying to view the details", Adminuser)
	// this is the username to which we want to get
	targetUsername := c.Request.PostFormValue("username")

	// Query the database to find the user based on the target username
	var user db.User
	if err := db.DBASE.Where("username = ?", targetUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the user data as a response
	c.JSON(http.StatusOK, user)

}

func DeleteUserData(c *gin.Context) {
	// admin user can delete the user details
	username, err := auth.ExtractUsernameFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	fmt.Println("user trying to delete the userdata", username)

	targetUsername := c.Request.PostFormValue("username")

	result := db.DBASE.Where("username = ?", targetUsername).Delete(&db.User{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user data"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User data deleted"})

}

func UpdateUserDate(c *gin.Context) {
	//admin user can update the user details
	username, err := auth.ExtractUsernameFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	fmt.Println("admin user trying to update user details", username)
	targetUsername := c.Request.PostFormValue("username")
	newUsername := c.Request.PostFormValue("newusername")

	result := db.DBASE.Model(&db.User{}).Where("username = ?", targetUsername).Update("username", newUsername)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user data"})
		return
	}

	// Check the affected rows to determine if the user data was updated
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User data updated"})
}
