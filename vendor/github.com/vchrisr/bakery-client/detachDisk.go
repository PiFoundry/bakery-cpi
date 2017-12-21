package bakeryclient

import (
	"fmt"
	"net/http"
)

func (c *Client) DetachDisk(piId, diskId string) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%v/oven/%v/disks/%v", c.url, piId, diskId), nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DetachDisk returned status code %v", resp.StatusCode)
	}

	return nil
}
