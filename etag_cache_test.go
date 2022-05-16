package etagcache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestAddEtag(t *testing.T) {
	// Create http mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	orders, err := ioutil.ReadFile("orderbook.json")
	if err != nil {
		t.Errorf("Error while fetching order json response. %v", err)
	}
	httpmock.RegisterResponder("GET", "https://api.kite.trade/orders",
		httpmock.NewStringResponder(200, string(orders)))

	// Create new cache instance
	cacheClient := New(ClientParam{DirPath: ""})
	url := "https://api.kite.trade/orders"

	req, _ := http.NewRequest("GET", url, nil)

	// HTTP request params
	reqParam := RequestParam{
		ReqClient: req,
		headers: map[string]string{
			"User-Agent":     "gokiteconnect/4.0.2",
			"x-kite-version": "3",
			"authorization":  "token api_key:access_token",
		},
	}
	// Store and fetch response/cache data
	response, err := cacheClient.HandleEtagCache(reqParam, url)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	var js map[string]interface{}
	json.Unmarshal([]byte(response), &js)
	assert.Equal(t, js["status"], "success", "Orderbook request failed.")
}
