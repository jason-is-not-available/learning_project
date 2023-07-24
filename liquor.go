// Writing a basic HTTP server is easy using the
// `net/http` package.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var mLiquors = map[string]int{"bourbon": 3, "vodka": 2}

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
		fail(w, "/ detected. Not acceptable")
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
			fail(w, "/ detected. Not acceptable")
			return
		}

		endpoint = endpoint[:len(endpoint)-1]
	}

	aEndpoint := strings.Split(endpoint, "/")

	// All valid use starts with /liquors
	if aEndpoint[0] != "liquors" {
		fail(w, "Another error")
		return
	}

	// Max valid
	if len(aEndpoint) > 2 {
		fail(w, "too many arguments")
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

	// switch aEndpoint[1] {
	// case "liquors":
	// 	fmt.Println("Case liquors")
	// 	return
	// }

}

func fail(w http.ResponseWriter, err string) {

	// fmt.Println("You fucked up. Do it right next time")

	w.WriteHeader(http.StatusInternalServerError)

	fmt.Fprint(w, err)

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

		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprint(w, err)
		return
	}

	/*
		uri
		https://stackoverflow.com/questions/31480710/validate-url-with-standard-package-in-go
	*/

	var liquorSlice []sLiquors

	//myLiquors := list.New()
	for key, value := range mLiquors {
		item := sLiquors{
			Type:   key,
			Amount: value,
		}
		liquorSlice = append(liquorSlice, item)
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

func liquorsType(w http.ResponseWriter, req *http.Request, liquorType string) {

	if req.Method != "GET" {

		err := "This only works with GET"

		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprint(w, err)
		return
	}

	inStock := mLiquors[liquorType]

	fmt.Println(inStock)

	returnInventory := sLiquors{
		Type:   liquorType,
		Amount: inStock,
	}

	fmt.Println("Print returnInventory", returnInventory)

	liquorsJson, _ := json.Marshal(returnInventory)

	fmt.Fprintf(w, string(liquorsJson))

	return

}

/*
This endpoint should receive a json object and add the amount to the
existing amount (or create the new entry). An example POST is below.
The response should be the corresponsing current total amount.*/

func addLiquors(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {

		err := "This only works with POST"

		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprint(w, err)
		return
	}

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	var addRequest sLiquors

	json.Unmarshal([]byte(body), &addRequest)

	mLiquors[addRequest.Type] += addRequest.Amount

	addRequest.Amount = mLiquors[addRequest.Type]

	liquorsJson, _ := json.Marshal(addRequest)
	fmt.Fprintf(w, string(liquorsJson))

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

		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprint(w, err)
		return
	}

	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()

	var removeRequest sLiquors

	json.Unmarshal([]byte(body), &removeRequest)

	quantity, exist := mLiquors[removeRequest.Type]

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

	mLiquors[removeRequest.Type] = mLiquors[removeRequest.Type] - removeRequest.Amount
	fmt.Println("Now we have", removeRequest.Type, mLiquors[removeRequest.Type])

	removeRequest.Amount = mLiquors[removeRequest.Type]

	liquorsJson, _ := json.Marshal(removeRequest)
	fmt.Fprintf(w, string(liquorsJson))

}

func main() {

	// Identify the endpoint and decide what to do with it
	// http.HandleFunc("/", myHandler)
	http.HandleFunc("/liquors/add", addLiquors)
	http.HandleFunc("/liquors/remove", removeLiquors)
	http.HandleFunc("/liquors", liquors)

	// Keep this around
	// http.HandleFunc("/liquors", liquors)

	// Finally, we call the `ListenAndServe` with the port
	// and a handler. `nil` tells it to use the default
	// router we've just set up.
	fmt.Println("Starting Server on 8090...")
	http.ListenAndServe(":8090", nil)

}
