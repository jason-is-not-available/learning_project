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

// var db *sql.DB = connectDB()

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

	fmt.Println("Create table if")
	query := "CREATE TABLE IF NOT EXISTS liquors (type varchar(255) not null UNIQUE, quantity int not null)"
	_, err := db.Query(query)
	if err != nil {
		fmt.Println("Error making table")
	}

	rows, err := db.Query("select * from liquors")
	if err != nil {
		fmt.Println("I dunno")
	}
	defer rows.Close()

	if rows.Next() {
		return
	}

	fmt.Println("Empty table. Populating")
	query = "INSERT INTO liquors (type, quantity) VALUES ('bourbon', 5), ('good bourbon', 7), ('vodka', 4), ('gin', 14), ('tequila', 5000)"
	_, err = db.Query(query)
	if err != nil {
		fmt.Println("Error populating table")
	}
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

func dbList(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		rows, err := db.Query("select * from liquors")

		if err != nil {
			fmt.Println("error")
		}
		defer rows.Close()

		var items []item
		var item item

		for rows.Next() {
			err = rows.Scan(&item.Type, &item.Amount)
			if err != nil {
				fmt.Println("error")
			}
			items = append(items, item)
		}
		c.JSON(http.StatusOK, items)
	}
	return gin.HandlerFunc(fn)
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
	// Is this check even needed?
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

func dbType(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		requestedType := "%" + strings.ToLower(c.Param("type")) + "%"

		rows, err := db.Query(`select * from liquors where type like $1`, requestedType)
		if err != nil {
			fmt.Println("Error querying for types")
		}
		defer rows.Close()

		var items []item
		var item item

		for rows.Next() {
			err = rows.Scan(&item.Type, &item.Amount)
			if err != nil {
				fmt.Println("error")
			}
			items = append(items, item)
		}

		if len(items) == 0 {
			item.Type = c.Param("type")
			item.Amount = 0
			c.JSON(http.StatusOK, item)
			return
		}

		c.JSON(http.StatusOK, items)

	}
	return gin.HandlerFunc(fn)
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

func dbAdd(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var add item
		if err := c.BindJSON(&add); err != nil {
			fail(c, "Add failed. Do it better.")
			return
		}

		rows, err := db.Query(`select * from liquors where type = $1`, add.Type)

		if err != nil {
			fmt.Println("Error selecting")
		}
		defer rows.Close()

		// var items []inStock
		var inStock item

		if rows.Next() {
			err = rows.Scan(&inStock.Type, &inStock.Amount)
			if err != nil {
				fmt.Println("error")
			}

			_, err = db.Query(`update liquors set quantity = $1 where type = $2`, add.Amount+inStock.Amount, add.Type)

			if err != nil {
				fmt.Println("error")
			}

			inStock.Amount += add.Amount
		} else {
			fmt.Println("Didn't have a next")

			_, err = db.Query(`INSERT INTO liquors (type, quantity) VALUES ($1, $2)`, add.Type, add.Amount)
			if err != nil {
				fmt.Println("error")
			}

			inStock.Type = add.Type
			inStock.Amount = add.Amount

		}

		c.JSON(http.StatusOK, inStock)

	}
	return gin.HandlerFunc(fn)
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

func dbRemove(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var remove item
		if err := c.BindJSON(&remove); err != nil {
			fail(c, "Fuck you. Do it better.")
			return
		}
		// Query for object
		rows, err := db.Query(`select * from liquors where type = $1`, remove.Type)
		if err != nil {
			if err == sql.ErrNoRows {
				// No result
				fail(c, "Not Enough Liquor")
			}
		}
		defer rows.Close()

		var inStock item

		if rows.Next() {
			err = rows.Scan(&inStock.Type, &inStock.Amount)
			if err != nil {
				fmt.Println("error")
			}
		}

		inStock.Amount = inStock.Amount - remove.Amount

		if inStock.Amount < 0 {
			fail(c, "Not Enough Liquor")
			return
		}

		_, err = db.Query(`update liquors set quantity = $1 where type = $2`, inStock.Amount, remove.Type)
		if err != nil {
			fmt.Println("error")
		}

		c.JSON(http.StatusOK, inStock)

	}
	return gin.HandlerFunc(fn)

}

func fail(c *gin.Context, err string) {

	c.String(http.StatusInternalServerError, err)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	db := connectDB()
	defer db.Close()

	populateTable(db)

	/*
	   Is this worth using?
	   https://bun.uptrace.dev/guide/complex-queries.html#parsing-request-params
	*/

	router := gin.Default()

	// router.GET("/liquors", inventoryList)
	router.GET("/liquors", dbList(db))

	// router.GET("/liquors/:type", inventoryType)
	router.GET("/liquors/:type", dbType(db))

	// router.POST("/liquors/add", inventoryAdd)
	router.POST("/liquors/add", dbAdd(db))

	// router.POST("/liquors/remove", inventoryRemove)
	router.POST("/liquors/remove", dbRemove(db))
	log.Fatal(router.Run(":8090"))
}
