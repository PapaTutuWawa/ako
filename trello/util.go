package trello

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	API_BASE_URL = "https://api.trello.com/1"
)

// Builds the url that can be used to query the trello API
func BuildRequestUrl(endpoint, key, token string) string {
	return API_BASE_URL + endpoint + "?key=" + key + "&token=" + token
}

// Performs a Get request and returns the data parsed by json.Unmarshal
func GetUnmarshalledData(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Read the Response Body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Let go parse the JSON
	var data interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
