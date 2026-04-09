package services

import (
	"api/models"
)

//var err error

func PopulateProductFamilies(productFamilies []models.ProductFamily) error {
	DB.Delete(&models.ProductFamily{}) //Clear out the db

	for _, family := range productFamilies {
		DB.Create(&family)
	}
	return nil
}
