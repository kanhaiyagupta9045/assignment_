package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/car_management/databases"
	"github.com/kanhaiyagupta9045/car_management/helpers"
	"github.com/kanhaiyagupta9045/car_management/models"
	"github.com/kanhaiyagupta9045/car_management/utils"
)

func AddCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var car models.Car

		if err := c.BindJSON(&car); err != nil {
			utils.ErrorResponse(c, err, http.StatusBadRequest)
			return
		}
		if validationErr := validate.Struct(&car); validationErr != nil {
			utils.ErrorResponse(c, validationErr, http.StatusBadRequest)
			return
		}
		user, err := helpers.GetUserFromCookie(c)
		if err != nil {
			utils.ErrorResponse(c, errors.New("error while retrieving user from cookie"), http.StatusInternalServerError)
			return
		}

		insertCarQuery := `INSERT INTO cars (user_id, car_name, tags, description, car_type, car_company, dealer) VALUES (?, ?, ?, ?, ?, ?, ?)`
		result, err := databases.DB.Exec(insertCarQuery, user.User_Id, car.CarName, car.Tags, car.Description, car.CarType, car.CarCompany, car.Dealer)
		if err != nil {
			utils.ErrorResponse(c, errors.New("failed to create car"), http.StatusInternalServerError)
			return
		}

		carID, err := result.LastInsertId()
		if err != nil {
			utils.ErrorResponse(c, errors.New("failed to retrieve car ID"), http.StatusInternalServerError)
			return
		}

		insertImageQuery := `INSERT INTO car_images (car_id, image_url) VALUES (?, ?)`
		fmt.Println(len(car.Images))
		for _, image := range car.Images {
			fmt.Printf("Inserting image for CarID %d: %s\n", carID, image.ImageURL)

			_, err := databases.DB.Exec(insertImageQuery, carID, image.ImageURL)
			if err != nil {
				utils.ErrorResponse(c, errors.New("failed to save car images"), http.StatusInternalServerError)
				return
			}

		}

		c.JSON(http.StatusCreated, gin.H{"message": "Car created successfully", "car_id": carID})

	}
}

func GetCarAddedByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := helpers.GetUserFromCookie(c)
		if err != nil {
			utils.ErrorResponse(c, errors.New("error retrieving user from cookie"), http.StatusInternalServerError)
			return
		}

		query := `SELECT car_id, user_id, car_name, tags, description, car_type, car_company, dealer FROM cars WHERE user_id = ?`
		rows, err := databases.DB.Query(query, user.User_Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cars", "message": err.Error()})
			return
		}
		defer rows.Close()

		var cars []models.Car
		for rows.Next() {
			var car models.Car
			if err := rows.Scan(&car.CarID, &car.UserID, &car.CarName, &car.Tags, &car.Description, &car.CarType, &car.CarCompany, &car.Dealer); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read car data", "message": err.Error()})
				return
			}

			imageQuery := `SELECT image_id, car_id, image_url FROM car_images WHERE car_id = ?`
			imageRows, err := databases.DB.Query(imageQuery, car.CarID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch car images", "message": err.Error()})
				return
			}
			defer imageRows.Close()

			var images []models.CarImage
			for imageRows.Next() {
				var carImage models.CarImage
				if err := imageRows.Scan(&carImage.ImageID, &carImage.CarID, &carImage.ImageURL); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read car image data", "message": err.Error()})
					return
				}
				images = append(images, carImage)
			}
			car.Images = images
			cars = append(cars, car)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during car retrieval", "message": err.Error()})
			return
		}

		if len(cars) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No cars found for this user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"cars": cars})
	}
}

func SearchCarsByKeyowrd() gin.HandlerFunc {
	return func(c *gin.Context) {

		keyword := c.DefaultQuery("keyword", "")

		if keyword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search keyword is required"})
			return
		}

		query := `
		SELECT car_id, car_name, tags, description, car_type, car_company, dealer
		FROM cars
		WHERE (car_name LIKE ? OR description LIKE ? OR tags LIKE ?)`

		searchPattern := "%" + keyword + "%"
		rows, err := databases.DB.Query(query, searchPattern, searchPattern, searchPattern)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cars", "message": err.Error()})
			return
		}
		defer rows.Close()

		var cars []models.Car
		for rows.Next() {
			var car models.Car
			if err := rows.Scan(&car.CarID, &car.CarName, &car.Tags, &car.Description, &car.CarType, &car.CarCompany, &car.Dealer); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read car data", "message": err.Error()})
				return
			}

			imageQuery := `SELECT image_url,car_id,image_id FROM car_images WHERE car_id = ?`
			imageRows, err := databases.DB.Query(imageQuery, car.CarID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch car images", "message": err.Error()})
				return
			}
			defer imageRows.Close()

			var images []models.CarImage
			for imageRows.Next() {
				var carimage models.CarImage
				if err := imageRows.Scan(&carimage.ImageURL, &carimage.CarID, &carimage.ImageID); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read car image data", "message": err.Error()})
					return
				}
				images = append(images, carimage)
			}

			car.Images = images

			cars = append(cars, car)
		}

		if len(cars) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No cars found matching the search criteria"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"cars": cars})
	}
}

func GetCarDetails() gin.HandlerFunc {
	return func(c *gin.Context) {

		carID := c.Param("car_id")
		query := `
		SELECT car_id, car_name, tags, description, car_type, car_company, dealer
		FROM cars
		WHERE car_id = ?`
		var car models.Car

		err := databases.DB.QueryRow(query, carID).Scan(
			&car.CarID,
			&car.CarName,
			&car.Tags,
			&car.Description,
			&car.CarType,
			&car.CarCompany,
			&car.Dealer,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch car details", "message": err.Error()})
			}
			return
		}

		imageQuery := `SELECT image_url,car_id,image_id FROM car_images WHERE car_id = ?`
		imageRows, err := databases.DB.Query(imageQuery, car.CarID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch car images", "message": err.Error()})
			return
		}
		defer imageRows.Close()
		var images []models.CarImage
		for imageRows.Next() {
			var image models.CarImage
			if err := imageRows.Scan(&image.ImageURL, &image.CarID, &image.ImageID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read car image data", "message": err.Error()})
				return
			}
			images = append(images, image)
		}
		car.Images = images
		c.JSON(http.StatusOK, gin.H{"car": car})
	}
}

func DeleteCar() gin.HandlerFunc {
	return func(c *gin.Context) {

		carID := c.Param("car_id")

		query := `
		DELETE FROM cars
		WHERE car_id = ?`

		result, err := databases.DB.Exec(query, carID)
		if err != nil {
			utils.ErrorResponse(c, err, http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			utils.ErrorResponse(c, err, http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Car Deleted Successfully"})
	}
}

func UpdateCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		carID := c.Param("car_id")
		var updatedCar models.Car

		if err := c.ShouldBindJSON(&updatedCar); err != nil {
			utils.ErrorResponse(c, err, http.StatusBadRequest)
			return
		}

		checkQuery := `SELECT car_name FROM cars WHERE car_id = ?`
		var carName string
		err := databases.DB.QueryRow(checkQuery, carID).Scan(&carName)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "No car exists with this car_id"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check car existence", "message": err.Error()})
			return
		}

		query := `
		UPDATE cars
		SET car_name = ?, tags = ?, description = ?, car_type = ?, car_company = ?, dealer = ?
		WHERE car_id = ?`

		_, err = databases.DB.Exec(query, updatedCar.CarName, updatedCar.Tags, updatedCar.Description,
			updatedCar.CarType, updatedCar.CarCompany, updatedCar.Dealer, carID)
		if err != nil {
			utils.ErrorResponse(c, err, http.StatusInternalServerError)
			return
		}

		if len(updatedCar.Images) > 0 {
			for _, image := range updatedCar.Images {

				existsQuery := `
					SELECT COUNT(1) 
					FROM car_images 
					WHERE car_id = ? AND image_url = ?`

				var exists int
				err := databases.DB.QueryRow(existsQuery, carID, image.ImageURL).Scan(&exists)
				if err != nil {
					utils.ErrorResponse(c, err, http.StatusInternalServerError)
					return
				}

				if exists == 0 {
					imageQuery := `
						INSERT INTO car_images (car_id, image_url)
						VALUES (?, ?)`

					_, err := databases.DB.Exec(imageQuery, carID, image.ImageURL)
					if err != nil {
						utils.ErrorResponse(c, err, http.StatusInternalServerError)
						return
					}
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Car updated successfully"})
	}
}
