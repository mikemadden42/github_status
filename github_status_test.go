package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestMain(m *testing.M) {
	color.NoColor = true
	os.Exit(m.Run())
}

func serveFixture(t *testing.T) *httptest.Server {
	t.Helper()
	data, err := os.ReadFile("testdata/summary.json")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(data)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestFetch(t *testing.T) {
	srv := serveFixture(t)

	s, err := fetch(srv.URL)
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if got := len(s.Components); got != 4 {
		t.Fatalf("got %d components, want 4", got)
	}
	if got := s.Status.Description; got != "Partial System Outage" {
		t.Errorf("status description = %q", got)
	}
}

func TestFetchErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    string
	}{
		{
			name: "server error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			want: "unexpected HTTP status: 500",
		},
		{
			name: "malformed JSON",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprint(w, "<html>not json</html>")
			},
			want: "decoding summary",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			_, err := fetch(srv.URL)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got error %v, want one containing %q", err, tt.want)
			}
		})
	}
}

func TestRender(t *testing.T) {
	s, err := fetch(serveFixture(t).URL)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := render(&buf, s); err != nil {
		t.Fatal(err)
	}

	want := fmt.Sprintf("Partial System Outage\n\n%-60s%s\n%-60s%s\n%-60s%s\n",
		"Git Operations", "operational",
		"Actions", "major_outage",
		"Pages", "degraded_performance")
	if got := buf.String(); got != want {
		t.Errorf("render output mismatch\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestStatusColor(t *testing.T) {
	tests := []struct {
		status string
		want   color.Attribute
	}{
		{"operational", color.FgGreen},
		{"degraded_performance", color.FgYellow},
		{"partial_outage", color.FgYellow},
		{"major_outage", color.FgRed},
		{"under_maintenance", color.FgBlue},
		{"something_new", color.FgYellow},
	}
	for _, tt := range tests {
		if !statusColor(tt.status).Equals(color.New(tt.want)) {
			t.Errorf("statusColor(%q) is not %v", tt.status, tt.want)
		}
	}
}
