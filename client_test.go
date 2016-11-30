package bioportal

import (
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		t.Skip("API_KEY not defined")
	}

	c := NewClient(apiKey)
	res, err := c.Search("audiology", nil)
	if err != nil {
		t.Fatal(err)
	}

	if res.Page == 0 {
		t.Fatal("no response")
	}
}
