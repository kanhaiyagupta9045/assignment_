package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kanhaiyagupta9045/car_management/databases"
	"github.com/kanhaiyagupta9045/car_management/helpers"
	"github.com/kanhaiyagupta9045/car_management/models"
	"github.com/kanhaiyagupta9045/car_management/utils"
)

var validate = validator.New()

func SignUpController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			utils.ErrorResponse(c, err, http.StatusBadRequest)
			return
		}
		if validationErr := validate.Struct(user); validationErr != nil {
			utils.ErrorResponse(c, validationErr, http.StatusBadRequest)
			return
		}

		checkemailquery := `SELECT email FROM users WHERE email = ?`
		var exitingEmail string
		err := databases.DB.QueryRow(checkemailquery, user.Email).Scan(&exitingEmail)

		if err == nil {
			utils.ErrorResponse(c, errors.New("email already exist"), http.StatusConflict)
			return
		}

		user.Password = helpers.HashPassword(user.Password)

		user.Created_at = time.Now()
		user.Updated_at = time.Now()

		insertuserquery := `INSERT INTO users(first_name,last_name,email,password,created_at,updated_at,deleted_at) VALUES(?,?,?,?,?,?,?)`
		result, err := databases.DB.Exec(insertuserquery, user.First_Name, user.Last_Name, user.Email, user.Password, user.Created_at, user.Updated_at, nil)
		if err != nil {
			utils.ErrorResponse(c, err, http.StatusInternalServerError)
		}
		id, err := result.LastInsertId()
		if err != nil {
			utils.ErrorResponse(c, err, http.StatusInternalServerError)
			return
		}
		user.User_Id = uint32(id)
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.User_Id})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var logindata models.LoginInfo

		if err := c.BindJSON(&logindata); err != nil {
			utils.ErrorResponse(c, err, http.StatusBadRequest)
			return
		}
		if validationErr := validate.Struct(logindata); validationErr != nil {
			utils.ErrorResponse(c, validationErr, http.StatusBadRequest)
			return
		}
		var existinguser models.User
		query := "SELECT user_id,first_name,last_name,email,password FROM users where email = ?"

		if err := databases.DB.QueryRow(query, logindata.Email).Scan(
			&existinguser.User_Id,
			&existinguser.First_Name,
			&existinguser.Last_Name,
			&existinguser.Email,
			&existinguser.Password,
		); err != nil {
			if err == sql.ErrNoRows {
				utils.ErrorResponse(c, errors.New("user not found"), http.StatusNotFound)
			} else {
				utils.ErrorResponse(c, err, http.StatusInternalServerError)
			}
			return
		}
		if ok := helpers.VerifyPassword(existinguser.Password, logindata.Password); !ok {
			utils.ErrorResponse(c, errors.New("provided password is incorrect"), http.StatusBadRequest)
			return
		}

		accesstoken, err := helpers.GenerateAccessToken(int(existinguser.User_Id))

		if err != nil {
			utils.ErrorResponse(c, err, http.StatusInternalServerError)
		}
		c.JSON(http.StatusOK, gin.H{"token": accesstoken})
	}
}
