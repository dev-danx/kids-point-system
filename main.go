package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type item struct {
	Name  string
	Point int32
}

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"content": "This is an index page...",
			"items":   readDatafileToStruct(),
		})
	})

	router.Run("localhost:8080")
}

func readDatafileToStruct() []item {
	// Open our jsonFile
	jsonFile, err := os.Open("data.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our array
	var result []item

	json.Unmarshal(byteValue, &result)
	return result
}
