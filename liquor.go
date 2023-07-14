// Writing a basic HTTP server is easy using the
// `net/http` package.
package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"
)

type response1 struct {
	Page   int
	Fruits []string
}

// Probably dont want a struct for this, but we'll start here
// might actually have use for it
type sLiquors struct {
	Type   []string
	Amount int
}

var mLiquors = map[string]int{"bourbon": 8, "vodka": 2}

// A fundamental concept in `net/http` servers is
// *handlers*. A handler is an object implementing the
// `http.Handler` interface. A common way to write
// a handler is by using the `http.HandlerFunc` adapter
// on functions with the appropriate signature.
func hello(w http.ResponseWriter, req *http.Request) {

	// Functions serving as handlers take a
	// `http.ResponseWriter` and a `http.Request` as
	// arguments. The response writer is used to fill in the
	// HTTP response. Here our simple response is just
	// "hello\n".

	page := &response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}

	//reqJson, _ := json.Marshal(req)
	//fmt.Fprintf(w, string(reqJson))

	fmt.Printf(req.Method)

	reqDump, _ := httputil.DumpRequest(req, true)
	fmt.Fprintf(w, string(reqDump))

	pageJson, _ := json.Marshal(page)
	fmt.Fprintf(w, string(pageJson))

	const menuLength int = 2
}

// Returns the entire map of liquors and their amounts
// in JSON format.
func liquors(w http.ResponseWriter, req *http.Request) {

	// make an array of structs from the map
	// then marshal that array

	// This should be a global int?
	// var arr [menuLength]sLiquors
	var arr [2]sLiquors

	// for i := 0; i < len(arr); i++ {
	// 	// ..cant iterate a map via index..
	// }

	for key, value := range mLiquors {
		i := 0
		// arr[i].Type = key
		// arr[i].Type = fmt.Sprintf(key)
		// Key is not a string and does not want to become one?
		// arr[i].Amount, _ = strconv.Atoi(value)
		arr[i].Amount = value
		fmt.Printf(key)
		// fmt.Printf(value)
		i++
	}

	myList := list.New()
	for key, value := range mLiquors {
		item := &sLiquors{
			Type:   []string{key},
			Amount: value,
		}
		fmt.Println(reflect.TypeOf(item))

		myList.PushBack(item)
	}

	one := myList.Front()
	fmt.Println(reflect.TypeOf(one))

	// var len = len(mLiquors)

	// var arr [len]sLiquors

	// for k, v := range mLiquors {
	// 	arr = append(arr, k, v)
	// }
	/*
		fmt.Println("Array ended up being: ")
		ftm.Println(arr) */

	pageJson, _ := json.Marshal(mLiquors)
	fmt.Fprintf(w, string(pageJson))

}

func main() {

	//fmt.Println("map:", mLiquors)

	// We register our handlers on server routes using the
	// `http.HandleFunc` convenience function. It sets up
	// the *default router* in the `net/http` package and
	// takes a function as an argument.
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/liquors", liquors)

	// Finally, we call the `ListenAndServe` with the port
	// and a handler. `nil` tells it to use the default
	// router we've just set up.
	http.ListenAndServe(":8090", nil)

}
