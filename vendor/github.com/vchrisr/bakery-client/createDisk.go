package bakeryclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) CreateDisk() (string, error) {
	req, _ := http.NewRequest("POST", c.url+"/disks", nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("CreateDisk returned status code %v", resp.StatusCode)
	}

	var disk Disk
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &disk)
	if err != nil {
		return "", err
	}

	return disk.ID, nil
}
