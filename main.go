package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type User struct {
	Email    string
	Password string
}

const (
	Db_USER     = "admin"
	Db_PASSWORD = "alypsok"
	Db_NAME     = "onelab5"
)

func main() {
	//connect to db
	DbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, Db_USER, Db_PASSWORD, Db_NAME)

	db, connectError := sql.Open("postgres", DbInfo)

	if connectError != nil {
		log.Fatal(connectError)
	}

	pingError := db.Ping()

	if pingError != nil {
		log.Fatal(pingError)
	}

	defer db.Close()

	// accounts := make(map[string]User)
	// simple user-register and login crud operations
	e := echo.New()

	//login user with email and password
	e.GET("/login/:email", func(c echo.Context) error {
		email := c.Param("email")
		password := c.QueryParam("password")

		var getUser User

		err := db.QueryRow("SELECT * FROM users WHERE email=$1 AND password=$2;", email, password).Scan(&getUser.Email, &getUser.Password)

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
