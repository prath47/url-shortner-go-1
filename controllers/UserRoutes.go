package controllers

import (
	"net/http"
	"os"
	"url-shortner-1/initializers"
	"url-shortner-1/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var reqBody struct {
		Name     string
		Email    string
		Password string
	}

	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Cannot Bind Json",
			},
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Cannot Hash Password",
			},
		})
		return
	}

	user := models.User{Name: reqBody.Name, Email: reqBody.Email, Password: string(hashedPassword)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Cannot Create User",
				"error":   result.Error.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User Created Successfully",
	})
}

func Login(c *gin.Context) {
	var reqBody struct {
		Email    string
		Password string
	}

	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid Params",
			},
		})
		return
	}

	// find user
	var user models.User
	initializers.DB.First(&user, "email = ?", reqBody.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Email Not found",
			},
		})
		return
	}

	// compare the password
	comparePassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))

	if comparePassword != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid Email or Password",
			},
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
	})

	SECRET := os.Getenv("SECRET")
	tokenString, _ := token.SignedString([]byte(SECRET))

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt", tokenString, 3600, "", "", false, true)
	c.JSON(200, gin.H{
		"user": user,
	})
}