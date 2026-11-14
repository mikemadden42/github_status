// Command github_status prints the status of each GitHub component.
package main

// https://tutorialedge.net/golang/consuming-restful-api-with-go/
// https://mholt.github.io/json-to-go/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type summary struct {
	Components []struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Showcase bool   `json:"showcase"`
	} `json:"components"`
	Status struct {
		Description string `json:"description"`
	} `json:"status"`
}

const summaryURL = "https://www.githubstatus.com/api/v2/summary.json"

func main() {
	check := flag.Bool("check", false, "exit with status 1 if any component is not operational")
	flag.Parse()

	s, err := fetch(summaryURL)
	checkErr(err)
	checkErr(render(os.Stdout, s))

	if *check && !operational(s) {
		os.Exit(1)
	}
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

// render writes the overall status followed by one line per component to w,
// with each component status colorized. Components not shown on the status
// page (showcase false) are skipped.
func render(w io.Writer, s summary) error {
	if _, err := fmt.Fprintf(w, "%s\n\n", s.Status.Description); err != nil {
		return err
	}
	for _, component := range s.Components {
		if !component.Showcase {
			continue
		}
		if _, err := fmt.Fprintf(w, "%-60s", component.Name); err != nil {
			return err
		}
		if _, err := statusColor(component.Status).Fprintln(w, component.Status); err != nil {
			return err
		}
	}
	return nil
}

// operational reports whether every component shown on the status page is
// operational.
func operational(s summary) bool {
	for _, component := range s.Components {
		if component.Showcase && component.Status != "operational" {
			return false
		}
	}
	return true
}

// statusColor returns the color used to display a component status.
func statusColor(status string) *color.Color {
	switch status {
	case "operational":
		return color.New(color.FgGreen)
	case "degraded_performance", "partial_outage":
		return color.New(color.FgYellow)
	case "major_outage":
		return color.New(color.FgRed)
	case "under_maintenance":
		return color.New(color.FgBlue)
	default:
		return color.New(color.FgYellow)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
