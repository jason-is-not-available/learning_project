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

func myHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("We got to my handler")

	endpoint := req.URL.Path
	method := req.Method

	fmt.Println("Method: ", method)

	// Lowercase and remove trailing "/" if it exists
	endpoint = strings.ToLower(endpoint)
	if strings.HasSuffix(endpoint, "/") {
		endpoint = endpoint[:len(endpoint)-1]
	}

	fmt.Println("testing slice", len(endpoint))
	fmt.Println(endpoint[:8])

	if len(endpoint) == 8 && endpoint[:8] == "/liquors" {
		fmt.Println("OMG LIQUOUROURSS")
	}

	switch endpoint {
	case "/liquors":
		fmt.Println("Case liquors")

	}

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

	println(req.URL.Path)

	whatis := req.URL.Path
	println(whatis)
	println(reflect.TypeOf(liquorSlice).String())

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
