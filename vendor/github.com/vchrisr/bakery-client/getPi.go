package bakeryclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetPi(piId string) (PiInfo, error) {
	resp, err := c.httpClient.Get(c.url + "/oven/" + piId)
	if err != nil {
		return PiInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PiInfo{}, nil
	}

	var parsedResponse PiInfo
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &parsedResponse)
	if err != nil {
		return PiInfo{}, err
	}

	return parsedResponse, nil
}
