package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type User struct {
	Email    string
	Password string
}

func main() {
	//connect to db
	var db *sqlx.DB
	var err error
	godotenv.Load(".env")

	db, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_NAME")))

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// simple user-register and login crud operations
	e := echo.New()

	//login user with email and password
	e.GET("/login/:email", func(c echo.Context) error {
		email := c.Param("email")
		password := c.QueryParam("password")

		var getUser User

		err := db.Get(&getUser, "SELECT * FROM users WHERE email=$1 AND password=$2;", email, password)

		if err == sql.ErrNoRows {
			return c.String(http.StatusBadRequest, "User does not exist")
		} else if err != nil {
			return c.String(http.StatusBadGateway, fmt.Sprintf("%v", err))
		} else {
			return c.String(http.StatusOK, "OK")
		}
	})

	// updating user password
	e.PATCH("/update/:email", func(c echo.Context) error {
		email := c.Param("email")
		oldPassword := c.QueryParam("password")
		newPassword := c.FormValue("password")

		rows, err := db.Exec("UPDATE users SET password=$1 WHERE email=$2 AND password=$3;", newPassword, email, oldPassword)

		if num, _ := rows.RowsAffected(); num == 0 {
			return c.String(http.StatusBadRequest, "Invalid user details")
		}
		if err == nil {
			return c.String(http.StatusOK, "OK")
		} else {
			return c.String(http.StatusBadGateway, fmt.Sprintf("%v", err))
		}
	})

	//creating user
	e.POST("/create-user", func(c echo.Context) error {
		newUser := new(User)

		bindingErr := c.Bind(newUser)

		if bindingErr != nil {
			return bindingErr
		}

		if len(newUser.Email) == 0 || len(newUser.Password) == 0 {
			return c.String(http.StatusBadRequest, "Invalid body")
		}

		_, err := db.Exec("INSERT INTO users (email, password) VALUES ($1, $2);", newUser.Email, newUser.Password)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "Error: User already exists")
		} else {
			return c.String(http.StatusOK, "OK")
		}
	})

	// deleting user
	e.DELETE("/delete/:email", func(c echo.Context) error {
		email := c.Param("email")
		password := c.QueryParam("password")

		rows, err := db.Exec("DELETE FROM users WHERE email=$1 AND password=$2;", email, password)

		if num, _ := rows.RowsAffected(); num == 0 {
			return c.String(http.StatusBadRequest, "Invalid user details")
		}
		if err != nil {
			return c.String(http.StatusBadGateway, fmt.Sprintf("%v", err))
		} else {
			return c.String(http.StatusOK, "OK")
		}
	})

	e.Logger.Fatal(e.Start(":3000"))
}
