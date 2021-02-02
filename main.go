package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	users = map[int]*user{}
)

var errorServerMessage string

func fillTheMapOnStartServer() {
	jsonFile, err := os.Open("users.json")
	defer jsonFile.Close()
	if err != nil {
		fmt.Print("File not found")
		errorServerMessage = "DB file is not loaded, please try again later"
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &users)
}

func getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid parameter")
	}
	if users[id] == nil {
		return c.JSON(http.StatusNotFound, "No matches")
	}
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid parameter")
	}
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	if users[id] == nil {
		return c.JSON(http.StatusNotFound, "No matches")
	}
	users[id].Name = u.Name
	updateJSONFile()
	return c.JSON(http.StatusOK, users[id])
}

func createUser(c echo.Context) error {
	newID := maxKey() + 1
	u := &user{
		ID: newID,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	updateJSONFile()
	return c.JSON(http.StatusCreated, u)
}

func deleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid parameter")
	}
	if users[id] == nil {
		return c.JSON(http.StatusNotFound, "No matches")
	}
	delete(users, id)
	updateJSONFile()
	return c.JSON(http.StatusNoContent, "Delete is complete")
}

func updateJSONFile() {
	dataOut, err := json.MarshalIndent(&users, "", "  ")
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}
	err = ioutil.WriteFile("users.json", dataOut, 0)
	if err != nil {
		log.Fatal("Cannot write updated settings file:", err)
	}
}

func getAllUsers(c echo.Context) error {
	if errorServerMessage != "" {
		return c.JSON(http.StatusInternalServerError, errorServerMessage)
	}
	return c.JSON(http.StatusOK, users)
}

func maxKey() int {
	var maxValueOfKey int
	for maxValueOfKey = range users {
		break
	}
	for n := range users {
		if n > maxValueOfKey {
			maxValueOfKey = n
		}
	}
	return maxValueOfKey
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fillTheMapOnStartServer()

	e.GET("/users", getAllUsers)
	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}
