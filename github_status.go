package main

// https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

type summary struct {
	Components []struct {
		CreatedAt          string      `json:"created_at"`
		Description        string      `json:"description"`
		Group              bool        `json:"group"`
		GroupID            interface{} `json:"group_id"`
		ID                 string      `json:"id"`
		Name               string      `json:"name"`
		OnlyShowIfDegraded bool        `json:"only_show_if_degraded"`
		PageID             string      `json:"page_id"`
		Position           int64       `json:"position"`
		Showcase           bool        `json:"showcase"`
		Status             string      `json:"status"`
		UpdatedAt          string      `json:"updated_at"`
	} `json:"components"`
	Incidents []interface{} `json:"incidents"`
	Page      struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		TimeZone  string `json:"time_zone"`
		UpdatedAt string `json:"updated_at"`
		URL       string `json:"url"`
	} `json:"page"`
	ScheduledMaintenances []interface{} `json:"scheduled_maintenances"`
	Status                struct {
		Description string `json:"description"`
		Indicator   string `json:"indicator"`
	} `json:"status"`
}

func main() {
	response, err := http.Get("https://kctbh9vrtdwd.statuspage.io/api/v2/summary.json")
	checkErr(err)

	responseData, err := ioutil.ReadAll(response.Body)
	checkErr(err)

	var responseObject summary
	err = json.Unmarshal(responseData, &responseObject)
	checkErr(err)

	for _, component := range responseObject.Components {
		fmt.Printf("%-60s", component.Name)
		if component.Status == "operational" {
			color.Set(color.FgGreen)
		} else if component.Status == "major_outage" {
			color.Set(color.FgRed)
		} else {
			color.Set(color.FgYellow)
		}
		fmt.Printf("%s\n", component.Status)
		color.Unset()
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
