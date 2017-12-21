package bakeryclient

import "net/http"

func (c *Client) DiskExists(diskId string) (bool, error) {
	resp, err := c.httpClient.Get(c.url + "/disks/" + diskId)
	if err == nil && resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, err
}
