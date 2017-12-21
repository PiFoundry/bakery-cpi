package bakeryclient

import (
	"fmt"
	"net/http"
)

func (c *Client) DeleteDisk(diskId string) error {
	req, _ := http.NewRequest("DELETE", c.url+"/disks/"+diskId, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DeleteDisk returned status code %v", resp.StatusCode)
	}

	return nil
}
