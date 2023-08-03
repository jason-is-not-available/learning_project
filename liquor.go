/*
Guide for gin
https://earthly.dev/blog/golang-gin-framework/


*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Probably dont want a struct for this, but we'll start here
// might actually have use for it
type item struct {
	Type   string
	Amount int
}

var inventoryMap = map[string]int{"bourbon": 3, "vodka": 2, "nice bourbon": 1}

/*
This looks like a mess.
Should I make more funtions that get called from here?
If so.. what should I move out of here?
Maybe a: string whatAreWeDoing(endpoint, method) function?
Then this would just be:
next := whatAreWeDoing(e, m)
switch next {
case: this

	liquors()

case: that

		postWhatever()
	}

..or is that basically what this is

Also, return isnt doing exactly what I want
Its returning to server.go, going through that, and then
starting over again at the top of this function with the
request url set to /favicon.ico

Then it runs down and hits a fail condition.
I could make logic to deal specifically with favicon.ico,
but that seems stupid. I don't get what its doing, or why.
Reproduce by triggering the fail on line 76
*/
func liquorsHandler(w http.ResponseWriter, req *http.Request) {

	if !(isGet(w, req)) {
		return
	}

	endpoint := req.URL.Path
	method := req.Method

	// please let me compile
	fmt.Println("Method: ", method)

	// handle /
	if len(endpoint) == 1 {
		oldFail(w, "/ detected. Not acceptable")
		return
	}

	endpoint = strings.ToLower(endpoint)

	// Deal with ////
	if strings.HasPrefix(endpoint, "/") {
		endpoint = endpoint[1:]
	} else {
		fmt.Println("I dont think this should ever happen.")
	}

	if strings.HasSuffix(endpoint, "/") {
		if len(endpoint) == 1 {
			oldFail(w, "/ detected. Not acceptable")
			return
		}

		endpoint = endpoint[:len(endpoint)-1]
	}

	aEndpoint := strings.Split(endpoint, "/")

	// All valid use starts with /liquors
	if aEndpoint[0] != "liquors" {
		oldFail(w, "All valid use starts with /liquors")
		return
	}

	// Max valid
	if len(aEndpoint) > 2 {
		oldFail(w, "too many arguments")
		return
	}

	if len(aEndpoint) == 1 && aEndpoint[0] == "liquors" {
		fmt.Println("We're going to /liquors")
		liquors(w, req)
		return
	}

	if aEndpoint[1] == "add" {
		fmt.Println("Going to add")
		// addLiquors(w, req, aEndpoint[1])
		return
	}

	if aEndpoint[1] == "remove" {
		fmt.Println("Going to remove")
		// removeLiquors(w, req, aEndpoint[1])
		return
	}

	// All thats left is /liquors/type
	liquorsType(w, req, aEndpoint[1])

}

func isGet(w http.ResponseWriter, req *http.Request) bool {
	if req.Method != "GET" {
		err := "This only works with GET"
		oldFail(w, err)
		return false
	}
	return true

}

func isNegative(w http.ResponseWriter, i int) bool {

	if math.Signbit(float64(i)) {
		err := "Don't be that guy"
		oldFail(w, err)
		return true
	}
	return false
}

func ginIsNegative(c *gin.Context, i int) bool {
	if math.Signbit(float64(i)) {
		err := "Don't be that guy"
		fail(c, err)
		return true
	}
	return false
}

func oldFail(w http.ResponseWriter, err string) {

	w.WriteHeader(http.StatusInternalServerError)

	fmt.Fprint(w, err)

}

func fail(c *gin.Context, err string) {

	c.String(http.StatusInternalServerError, err)

}

// Returns the entire map of liquors and their quantities
// in JSON format.

/*
This one will now handle both /liquors and /liquors/TYPE
Take some logic from myhander.
Or use myhandler to decide what happens here
*/
func liquors(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		err := "This only works with GET"
		oldFail(w, err)
		return
	}

	/*
		uri
		https://stackoverflow.com/questions/31480710/validate-url-with-standard-package-in-go
		https://pkg.go.dev/net/url

	*/

	/*

	   call myhandler here?
	   No. call that from main. Then that calls this

	*/

	var liquorSlice []item

	//myLiquors := list.New()
	for key, value := range inventoryMap {
		item := item{
			Type:   key,
			Amount: value,
		}
		liquorSlice = append(liquorSlice, item)
	}

	liquorsJson, _ := json.Marshal(liquorSlice)

	fmt.Fprint(w, string(liquorsJson))
}

func liquorsType(w http.ResponseWriter, req *http.Request, liquorType string) {

	// Might also want to check for an empty string here?
	if req.Method != "GET" {

		err := "This only works with GET"
		oldFail(w, err)
		return
	}

	inStock := inventoryMap[liquorType]

	fmt.Println(inStock)

	returnInventory := item{
		Type:   liquorType,
		Amount: inStock,
	}

	fmt.Println("Print returnInventory", returnInventory)

	liquorsJson, _ := json.Marshal(returnInventory)

	fmt.Fprint(w, string(liquorsJson))

}

/*
This endpoint should receive a json object and add the amount to the
existing amount (or create the new entry). An example POST is below.
The response should be the corresponsing current total amount.*/

func addLiquors(w http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		err := "This only works with POST"
		oldFail(w, err)
		return
	}

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	var addRequest item

	json.Unmarshal([]byte(body), &addRequest)

	if isNegative(w, addRequest.Amount) {
		return
	}

	inventoryMap[addRequest.Type] += addRequest.Amount

	addRequest.Amount = inventoryMap[addRequest.Type]

	liquorsJson, _ := json.Marshal(addRequest)
	fmt.Fprint(w, string(liquorsJson))
}

/*
This endpoint should receive a json object and remove the amount from the existing amount.
If the number requested is more than the current total, a 500 error should be thrown.

An example POST is below. The response should be the corresponsing current total amount.

```
POST /liquors/remove {"type": "bourbon", "amount": 4}

Response (200):

{"type": "bourbon", "amount": 4}

OR

Response (500):

"Not Enough Liquor"
*/

func removeLiquors(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Made it to removeLiquors")
	if req.Method != "POST" {

		err := "This only works with POST"
		oldFail(w, err)
		return
	}

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	var removeRequest item

	json.Unmarshal([]byte(body), &removeRequest)

	if isNegative(w, removeRequest.Amount) {
		return
	}

	// // check for negatives
	// if math.Signbit(float64(removeRequest.Amount)) {
	// 	err := "Don't be that guy"
	// 	fail(w, err)
	// 	return
	// }

	quantity, exist := inventoryMap[removeRequest.Type]

	if !exist || (quantity < removeRequest.Amount) {
		// fmt.Println("We don't have any", removeRequest.Type)
		// return 500 + message

		w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("Not Enough Liquor"))

		fmt.Fprint(w, string("Not Enough Liquor"))
		return

	}

	/*
		It does exist and we do have enough
	*/

	inventoryMap[removeRequest.Type] = inventoryMap[removeRequest.Type] - removeRequest.Amount
	fmt.Println("Now we have", removeRequest.Type, inventoryMap[removeRequest.Type])

	removeRequest.Amount = inventoryMap[removeRequest.Type]

	liquorsJson, _ := json.Marshal(removeRequest)
	fmt.Fprint(w, string(liquorsJson))
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
	// c.IndentedJSON(http.StatusOK, inventorySlice)

}

func inventoryType(c *gin.Context) {

	// c.String(http.StatusOK, "Made it to ginLiquorsType")

	// var item item
	var items []item

	request := c.Param("type")

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

	/*
		var item item
		item.Type = c.Param("type")
		item.Amount = inventoryMap[item.Type]
		c.JSON(http.StatusOK, item)
	*/

	/*
		inventory, inStock := inventoryMap[item.Type]

		if inStock {
			// Reply with type and stock
			item.Amount = inventory
			c.JSON(http.StatusOK, item)
			return
		}

		s := "We dont, and never did have any: " + item.Type
		c.JSON(200, s)
	*/

}

func inventoryAdd(c *gin.Context) {

	// c.String(http.StatusOK, "Made it to addInventory")
	var add item

	if err := c.BindJSON(&add); err != nil {
		return
	}

	if ginIsNegative(c, add.Amount) {
		return
	}

	inventoryMap[add.Type] += add.Amount

	add.Amount = inventoryMap[add.Type]

	c.JSON(http.StatusOK, add)
}

func inventoryRemove(c *gin.Context) {

	// c.String(http.StatusOK, "Made it to removeInventory")

	var remove item

	if err := c.BindJSON(&remove); err != nil {
		return
	}

	if ginIsNegative(c, remove.Amount) {
		return
	}

	quantity, exist := inventoryMap[remove.Type]

	if !exist || (quantity < remove.Amount) {
		fail(c, "Not Enough Liquor")
		return
	}

	inventoryMap[remove.Type] = inventoryMap[remove.Type] - remove.Amount
	// fmt.Println("Now we have", remove.Type, inventoryMap[remove.Type])

	remove.Amount = inventoryMap[remove.Type]

	// liquorsJson, _ := json.Marshal(remove)
	// fmt.Fprint(w, string(liquorsJson))
	c.JSON(http.StatusOK, remove)
}

func main() {

	router := gin.Default()
	router.GET("/liquors", inventoryList)
	router.GET("/liquors/:type", inventoryType)
	router.POST("/liquors/add", inventoryAdd)
	router.POST("/liquors/remove", inventoryRemove)
	log.Fatal(router.Run(":8090"))

	/*
		https: //chenyitian.gitbooks.io/gin-web-framework/content/docs/8.html
		https://earthly.dev/blog/golang-gin-framework/
		https://stackoverflow.com/questions/48010954/json-response-in-golang-s-gin-returning-as-scrambled-data
		https://chenyitian.gitbooks.io/gin-web-framework/content/docs/39.html
	*/

	/*
		http.HandleFunc("/liquors/add", addLiquors)
		http.HandleFunc("/liquors/remove", removeLiquors)
		http.HandleFunc("/liquors/", liquorsHandler)
		http.HandleFunc("/liquors", liquorsHandler)
		fmt.Println("Starting Server on 8090...")
		http.ListenAndServe(":8090", nil)
	*/

	/*
		https://github.com/gin-gonic/gin#runnint-gin
	*/

}
