package initializers

import "url-shortner-1/models"

func MigrateDB() {
	DB.AutoMigrate(&models.User{}, &models.URL{})
}
