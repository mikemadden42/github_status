// Command github_status prints the status of each GitHub component.
package main

// https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get("https://kctbh9vrtdwd.statuspage.io/api/v2/summary.json")
	checkErr(err)
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		checkErr(fmt.Errorf("unexpected HTTP status: %s", response.Status))
	}

	var responseObject summary
	err = json.NewDecoder(response.Body).Decode(&responseObject)
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
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
