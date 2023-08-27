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

func createTable(db *sql.DB) {

	query := "CREATE TABLE IF NOT EXISTS liquors (type varchar(255) not null UNIQUE, quantity int not null)"
	_, err := db.Query(query)
	if err != nil {
		fmt.Println("Error making table")
	}

}

func populateTable(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		rows, err := db.Query("SELECT * FROM liquors")
		if err != nil {
			fmt.Println("I dunno")
		}
		defer rows.Close()

		if rows.Next() {
			c.JSON(http.StatusPreconditionFailed, "Table not empty")
			return
		}

		fmt.Println("Empty table. Populating")
		query := "INSERT INTO liquors (type, quantity) VALUES ('bourbon', 5), ('good bourbon', 7), ('vodka', 4), ('gin', 14), ('tequila', 5000)"
		_, err = db.Query(query)
		if err != nil {
			fmt.Println("Error populating table")
		}
		c.JSON(http.StatusOK, "Table populated")
	}
	return gin.HandlerFunc(fn)

}

func inventoryList(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM liquors")

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

func inventoryType(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		requestedType := "%" + strings.ToLower(c.Param("type")) + "%"

		rows, err := db.Query(`SELECT * FROM liquors WHERE type like $1`, requestedType)
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

func inventoryAdd(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var add item
		if err := c.BindJSON(&add); err != nil {
			fail(c, "Add failed. Do it better.")
			return
		}

		rows, err := db.Query(`SELECT * FROM liquors WHERE type = $1`, add.Type)

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

			_, err = db.Query(`UPDATE liquors SET quantity = $1 where type = $2`, add.Amount+inStock.Amount, add.Type)

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

func inventoryRemove(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var remove item
		if err := c.BindJSON(&remove); err != nil {
			fail(c, "Fuck you. Do it better.")
			return
		}
		// Query for object
		rows, err := db.Query(`SELECT * FROM liquors WHERE type = $1`, remove.Type)
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

		_, err = db.Query(`UPDATE liquors SET quantity = $1 WHERE type = $2`, inStock.Amount, remove.Type)
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

	createTable(db)
	populateTable(db)

	/*
	   Is this worth using?
	   https://bun.uptrace.dev/guide/complex-queries.html#parsing-request-params
	*/

	router := gin.Default()
	router.GET("/liquors", inventoryList(db))
	router.GET("/liquors/:type", inventoryType(db))
	router.POST("/liquors/add", inventoryAdd(db))
	router.POST("/liquors/remove", inventoryRemove(db))
	router.POST("/populate", populateTable(db))
	log.Fatal(router.Run(":8090"))
}
