// Writing a basic HTTP server is easy using the
// `net/http` package.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// Probably dont want a struct for this, but we'll start here
// might actually have use for it
type sLiquors struct {
	Type   string
	Amount int
}

var mLiquors = map[string]int{"bourbon": 8, "vodka": 2}

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
func myHandler(w http.ResponseWriter, req *http.Request) {
	// fmt.Println("We got to my handler")

	endpoint := req.URL.Path
	method := req.Method

	// please let me compile
	fmt.Println("Method: ", method)

	// handle /
	if len(endpoint) == 1 {
		fail(w)
		return
	}

	// Lowercase and remove bookend "/" if present
	endpoint = strings.ToLower(endpoint)
	if strings.HasPrefix(endpoint, "/") {
		endpoint = endpoint[1:]
	} else {
		fmt.Println("I dont think this should ever happen.")
	}

	if strings.HasSuffix(endpoint, "/") {
		if len(endpoint) == 1 {
			fail(w)
			return
		}

		endpoint = endpoint[:len(endpoint)-1]
	}

	aEndpoint := strings.Split(endpoint, "/")

	// All valid use starts with /liquors
	if aEndpoint[0] != "liquors" {
		fail(w)
		return
	}

	// Max valid
	if len(aEndpoint) > 2 {
		fail(w)
		return
	}

	if len(aEndpoint) == 1 {
		fmt.Println("We're going to /liquors")
		liquors(w, req)
		return
	}

	fmt.Println("we have two parameters")
	fmt.Println("Lets see about a map match")

	inStock := mLiquors[aEndpoint[1]]
	fmt.Println("From map", inStock)

	returnInventory := sLiquors{
		Type:   aEndpoint[1],
		Amount: inStock,
	}

	fmt.Println("Print returnInventory", returnInventory)

	liquorsJson, _ := json.Marshal(returnInventory)

	fmt.Fprintf(w, string(liquorsJson))
	return

	// if liquourRequest {
	// 	fmt.Println("Need to look this up")
	// } else {
	// 	fmt.Println("Tell them we have zero")
	// }

	// switch aEndpoint[1] {
	// case "liquors":
	// 	fmt.Println("Case liquors")
	// 	return
	// }

}

func fail(w http.ResponseWriter) {

	fmt.Println("You fucked up. Try again")
	derp := "You fucked up. Try again"

	// derpJson, _ := json.Marshal(derp)

	fmt.Fprintf(w, derp)

}

// Returns the entire map of liquors and their quantities
// in JSON format.
func liquors(w http.ResponseWriter, req *http.Request) {

	var liquorSlice []sLiquors

	//myLiquors := list.New()
	for key, value := range mLiquors {
		item := sLiquors{
			Type:   key,
			Amount: value,
		}
		liquorSlice = append(liquorSlice, item)
		//fmt.Println(item)
		//testJson, _ := json.Marshal(item)
		//fmt.Println(string(testJson))
		//myLiquors.PushBack(item)
	}

	// This seems to do noting
	fmt.Println(liquorSlice)

	fmt.Println(req.URL.Path)

	whatis := req.URL.Path
	fmt.Println(whatis)
	fmt.Println(reflect.TypeOf(liquorSlice).String())

	liquorsJson, _ := json.Marshal(liquorSlice)

	fmt.Fprintf(w, string(liquorsJson))

}

func main() {

	// Identify the endpoint and decide what to do with it

	http.HandleFunc("/", myHandler)

	// Keep this around
	// http.HandleFunc("/liquors", liquors)

	// Finally, we call the `ListenAndServe` with the port
	// and a handler. `nil` tells it to use the default
	// router we've just set up.
	fmt.Println("Starting Server on 8090...")
	http.ListenAndServe(":8090", nil)

}
