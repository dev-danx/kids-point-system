package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type item struct {
	Id    int32
	Name  string
	Point int32
}

type update struct {
	Id            string
	PointToRemove int32
}

func (u update) IdAsInt() int {
	i, err := strconv.Atoi(u.Id)
	if err != nil {
		panic(err)
	}
	return i
}

//go:embed templates/*
var f embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		//log.Fatal("Error loading .env file")
	}

	templ := template.Must(template.New("").ParseFS(f, "templates/*.tmpl"))

	router := gin.Default()
	router.SetHTMLTemplate(templ)
	router.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
			"items": readDatafileToStruct(),
		})
	})

	router.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.tmpl", gin.H{
			"items": readDatafileToStruct(),
		})
	})

	router.POST("/update", func(c *gin.Context) {
		fmt.Println("Update Received")
		var updateItem update
		if err := c.BindJSON(&updateItem); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(updateItem)
		items := readDatafileToStruct()
		for i, el := range items {
			if el.Id == int32(updateItem.IdAsInt()) {
				fmt.Println("Found Match!")
				fmt.Println("Points to remove", updateItem.PointToRemove)
				fmt.Println("Points Before remove", el.Point)
				newPoints := el.Point + updateItem.PointToRemove
				items[i].Point = newPoints
				fmt.Println(items[i].Point)

			}
		}
		updateData(items)
		return
	})

	fmt.Println("Server startet!")
	listenOn := os.Getenv("LISTENON")
	router.Run(listenOn)
}

func readDatafileToStruct() []item {
	// Open our jsonFile
	jsonFile, err := os.Open("data.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened data.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our array
	var result []item

	json.Unmarshal(byteValue, &result)
	return result
}

func updateData(items []item) {
	j, err := json.Marshal(items)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		f, _ := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		defer f.Close()

		f.WriteString(string(j))
		fmt.Println("Updated json file", string(j))
	}

}
