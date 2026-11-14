// Command github_status prints the status of each GitHub component.
package main

// https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"fmt"
	"io"
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

const summaryURL = "https://kctbh9vrtdwd.statuspage.io/api/v2/summary.json"

func main() {
	s, err := fetch(summaryURL)
	checkErr(err)
	checkErr(render(os.Stdout, s))
}

// fetch retrieves and decodes the status page summary at url.
func fetch(url string) (summary, error) {
	var s summary

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return s, err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return s, fmt.Errorf("unexpected HTTP status: %s", response.Status)
	}

	if err := json.NewDecoder(response.Body).Decode(&s); err != nil {
		return s, fmt.Errorf("decoding summary: %w", err)
	}
	return s, nil
}

// render writes one line per component to w, with the status colorized.
func render(w io.Writer, s summary) error {
	for _, component := range s.Components {
		c := color.New(color.FgYellow)
		if component.Status == "operational" {
			c = color.New(color.FgGreen)
		} else if component.Status == "major_outage" {
			c = color.New(color.FgRed)
		}
		if _, err := fmt.Fprintf(w, "%-60s", component.Name); err != nil {
			return err
		}
		if _, err := c.Fprintln(w, component.Status); err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
