package controllers

import (
	"net/http"
	"url-shortner-1/helpers"
	"url-shortner-1/initializers"
	"url-shortner-1/models"

	"github.com/gin-gonic/gin"
)

// encrypt the URL
func EncryptTheUrl(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(500, gin.H{"error": "user not found"})
		return
	}

	var ReqBody struct {
		Url string
	}
	if c.BindJSON(&ReqBody) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Cannot Bind Json",
			},
		})
		return
	}

	if len(ReqBody.Url) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "URL Cannot Be Empty",
			},
		})
		return
	}

	uniqueId := helpers.RandString(6)

	shotedUrlObject := models.URL{
		ShortUrl:   uniqueId,
		MainUrl:    ReqBody.Url,
		Fk_id_user: int(userID.(uint)),
	}

	initializers.DB.Create(&shotedUrlObject)

	c.JSON(200, shotedUrlObject)
}

// get the url
func RedirectToPage(c *gin.Context) {
	uniqueCode := c.Param("Url")

	if uniqueCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Url",
		})
		return
	}

	var Url models.URL
	// result := initializers.DB.Where("short_url = ?", uniqueCode).First(&Url)
	result := initializers.DB.Where(models.URL{
		ShortUrl: uniqueCode,
	}).First(&Url)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"redirectTo": Url.MainUrl,
	})
}

// get all user url
func GetAllUserUrl(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(500, gin.H{"error": "user not found"})
		return
	}

	var allShortedUrl []models.URL
	initializers.DB.Where("Fk_id_user = ?", userID).Find(&allShortedUrl)

	len := len(allShortedUrl)

	c.JSON(200, gin.H{
		"data":  allShortedUrl,
		"count": len,
	})
}

// delete the url
func DeleteUrl(c *gin.Context) {
	var reqBody struct{
		Url string
	}

	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Url not valid",
		})
		return
	}

	var Url models.URL
	result := initializers.DB.Where(models.URL{
		ShortUrl: reqBody.Url,
	}).First(&Url)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid Url",
		})
		print(result.Error)
		return
	}

	if initializers.DB.Delete(&Url, Url.ID).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error Deleting URL Url",
		})
		print(result.Error)
		return
	}

	c.JSON(200, gin.H{
		"message": "Record Delete Successfully",
	})
}
