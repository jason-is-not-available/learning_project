package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type item struct {
	Type   string `json:"type" binding:"required,lowercase,min=1"`
	Amount int    `json:"amount" binding:"required,min=0"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "lrsql_user"
	password = "swordfish"
	dbname   = "lrsql_pg"
)

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

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func connectDB() *sql.DB {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	return db
}

func populateTable(db *sql.DB) {
	query := "select * from liquors"
	if _, err := db.Query(query); err == nil {
		fmt.Println("Table appears to exist")
		// return

	} else {

		fmt.Println("Table does not exist. Creating")

		query = "CREATE TABLE IF NOT EXISTS liquors (type varchar(255) not null UNIQUE, quantity int not null)"
		_, err := db.Query(query)

		if err != nil {
			fmt.Println("Error making table")
		}

		query = "INSERT INTO liquors (type, quantity) VALUES ('bourbon', 5), ('vodka', 4), ('gin', 14), ('tequila', 5000)"
		_, err = db.Query(query)
		if err != nil {
			fmt.Println("Error populating table")
		}
	}

}

func main() {

	db := connectDB()
	defer db.Close()

	populateTable(db)

	router := gin.Default()
	router.GET("/liquors", inventoryList)
	router.GET("/liquors/:type", inventoryType)
	router.POST("/liquors/add", inventoryAdd)
	router.POST("/liquors/remove", inventoryRemove)
	log.Fatal(router.Run(":8090"))
}
