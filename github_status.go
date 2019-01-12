package main

// https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type messages []struct {
	Status    string    `json:"status"`
	Body      string    `json:"body"`
	CreatedOn time.Time `json:"created_on"`
}

func main() {
	response, err := http.Get("https://status.github.com/api/messages.json")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject messages
	json.Unmarshal(responseData, &responseObject)

	for _, message := range responseObject {
		fmt.Printf("%s - %s\n%s\n\n",
			message.CreatedOn,
			message.Status,
			message.Body)
	}
}
