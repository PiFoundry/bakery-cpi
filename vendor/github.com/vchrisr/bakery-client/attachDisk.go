package bakeryclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) AttachDisk(piId, diskId string) error {
	var associateRequest struct {
		DiskId string `json:"diskId"`
	}

	associateRequest.DiskId = diskId
	jsonBytes, _ := json.Marshal(associateRequest)
	resp, err := c.httpClient.Post(fmt.Sprintf("%v/oven/%v/disks", c.url, piId), "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AttachDisk returned status code %v", resp.StatusCode)
	}

	return nil
}
