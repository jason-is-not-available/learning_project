package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type item struct {
	Type   string `json:"type" binding:"required,min=1"`
	Amount int    `json:"amount" binding:"required,min=0"`
}

var inventoryMap = map[string]int{
	"bourbon":      3,
	"vodka":        2,
	"nice bourbon": 1,
}

func inventoryList(c *gin.Context) {

	var inventorySlice []item

	for key, value := range inventoryMap {
		item := item{
			Type:   key,
			Amount: value,
		}
		inventorySlice = append(inventorySlice, item)
	}

	c.JSON(http.StatusOK, inventorySlice)
}

func inventoryType(c *gin.Context) {

	request := strings.ToLower(c.Param("type"))
	var items []item

	for key, value := range inventoryMap {
		if strings.Contains(key, request) {
			item := item{
				Type:   key,
				Amount: value,
			}
			items = append(items, item)
		}
	}

	// Should this be an if / else? or an if with return
	if len(items) == 0 {
		item := item{
			Type:   request,
			Amount: 0,
		}
		c.JSON(http.StatusOK, item)
	} else {
		c.JSON(http.StatusOK, items)
	}
}

func inventoryAdd(c *gin.Context) {

	var add item

	if err := c.BindJSON(&add); err != nil {
		fail(c, "Fuck you. Do it better.")
		return
	}



	add.Type = strings.ToLower(add.Type)


	inventoryMap[add.Type] += add.Amount
	add.Amount = inventoryMap[add.Type]

	c.JSON(http.StatusOK, add)
}

func inventoryRemove(c *gin.Context) {

	var remove item

	if err := c.BindJSON(&remove); err != nil {
		fail(c, "Fuck you. Do it better.")
		return
	}


  

	remove.Type = strings.ToLower(remove.Type)


  quantity, exist := inventoryMap[remove.Type]

	//Inventory does not exist
	if !exist || (quantity < remove.Amount) {
		fail(c, "Not Enough Liquor")
		return
	}

	inventoryMap[remove.Type] = inventoryMap[remove.Type] - remove.Amount
	remove.Amount = inventoryMap[remove.Type]

	c.JSON(http.StatusOK, remove)
}

func fail(c *gin.Context, err string) {

	c.String(http.StatusInternalServerError, err)
}

func main() {

	router := gin.Default()
	router.GET("/liquors", inventoryList)
	router.GET("/liquors/:type", inventoryType)
	router.POST("/liquors/add", inventoryAdd)
	router.POST("/liquors/remove", inventoryRemove)
	log.Fatal(router.Run(":8090"))
}
