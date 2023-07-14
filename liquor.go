// Writing a basic HTTP server is easy using the
// `net/http` package.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type response1 struct {
	Page   int
	Fruits []string
}

// Probably dont want a struct for this, but we'll start here
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

	reqJson, _ := json.Marshal(req)
	fmt.Fprintf(w, string(reqJson))

	fmt.Printf(req.Method)

	reqDump, _ := httputil.DumpRequest(req, true)
	fmt.Fprintf(w, string(reqDump))

	pageJson, _ := json.Marshal(page)
	fmt.Fprintf(w, string(pageJson))

}

func main() {

	//fmt.Println("map:", mLiquors)

	// We register our handlers on server routes using the
	// `http.HandleFunc` convenience function. It sets up
	// the *default router* in the `net/http` package and
	// takes a function as an argument.
	http.HandleFunc("/hello", hello)

	// Finally, we call the `ListenAndServe` with the port
	// and a handler. `nil` tells it to use the default
	// router we've just set up.
	http.ListenAndServe(":8090", nil)

}
