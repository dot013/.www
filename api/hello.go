package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
)

type helloObj struct {
	Language string
	Hello    string
}

func getHelloList() ([]helloObj, error) {
	res, err := http.Get("https://raw.githubusercontent.com/novellac/multilanguage-hello-json/master/hello.json")
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var hellos []helloObj
	err = json.Unmarshal(bytes, &hellos)
	if err != nil {
		return nil, err
	}

	return hellos, nil
}

func Hello(w http.ResponseWriter, r *http.Request) {
	hellos, err := getHelloList()
	var hello string
	if err != nil {
		hello = "Welcome!"
	} else {
		hello = hellos[rand.IntN(len(hellos)-1)].Hello
	}

	fmt.Fprint(w, hello)
}
