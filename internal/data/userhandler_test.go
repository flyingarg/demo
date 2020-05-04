package data

import (
	"encoding/json"
	"testing"
)

func TestGetUser(t *testing.T) {
	r := Request{
		ID:     "",
		Items:  []Item{{
			ID:           "",
			Name:         "Tomato Soup",
			Manufacturer: "Nestle",
			Brand:        "Knor Soups",
			Category:     "Ready to Eat Foods",
			Images:       []Image{{
				URL:    "http://127.0.0.1:7890/test.jpg",
				ID:     "",
				Status: false,
				Error:  "",
			}},
			Status:       false,
		}},
		Status: false,
	}
	d, _ := json.Marshal(r)
	t.Logf("%s", d)
}