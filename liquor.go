// Writing a basic HTTP server is easy using the
// `net/http` package.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)


// Probably dont want a struct for this, but we'll start here
// might actually have use for it
type sLiquors struct {
	Type   string
	Amount int
}

var mLiquors = map[string]int{"bourbon": 8, "vodka": 2}

// Returns the entire map of liquors and their amounts
// in JSON format.
func liquors(w http.ResponseWriter, req *http.Request) {

	var liquorArray []sLiquors

	//myLiquors := list.New()
	for key, value := range mLiquors {
		item := sLiquors{
			Type:   key,
			Amount: value,
		}
		liquorArray = append(liquorArray, item)
		//fmt.Println(item)
		//testJson, _ := json.Marshal(item)
		//fmt.Println(string(testJson))
		//myLiquors.PushBack(item)
	}

	fmt.Println(liquorArray)

	liquorsJson, _ := json.Marshal(liquorArray)
	//fmt.Println(string(liquorsJson))
	fmt.Fprintf(w, string(liquorsJson))

}

func main() {

	http.HandleFunc("/liquors", liquors)

	// Finally, we call the `ListenAndServe` with the port
	// and a handler. `nil` tells it to use the default
	// router we've just set up.
	fmt.Println("Starting Server on 8090...")
	http.ListenAndServe(":8090", nil)

}
