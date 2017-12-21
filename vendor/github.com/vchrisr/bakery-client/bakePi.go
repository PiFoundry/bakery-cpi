package bakeryclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type bakeRequest struct {
	BakeformName string `json:"bakeformName"`
}

func (c *Client) BakePi(bakeformId string) (PiInfo, error) {
	brq := bakeRequest{
		BakeformName: bakeformId,
	}

	body, err := json.Marshal(brq)
	if err != nil {
		return PiInfo{}, err
	}
	resp, err := c.httpClient.Post(c.url+"/fridge", "application/json", bytes.NewReader(body))
	if err != nil {
		return PiInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return PiInfo{}, fmt.Errorf("BakePi returned status code %v", resp.StatusCode)
	}

	var pi PiInfo
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &pi)
	if err != nil {
		return PiInfo{}, err
	}

	return pi, nil
}
