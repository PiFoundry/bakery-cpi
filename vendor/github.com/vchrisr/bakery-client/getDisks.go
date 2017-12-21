package bakeryclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetDisks() ([]string, error) {
	var diskResponse struct {
		Disks map[string]Disk `json:"disks"`
	}

	resp, err := c.httpClient.Get(c.url + "/disks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetDisks returned status code %v", resp.StatusCode)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &diskResponse)
	if err != nil {
		return nil, err
	}

	diskIds := make([]string, len(diskResponse.Disks))
	i := 0
	for key := range diskResponse.Disks {
		diskIds[i] = key
		i++
	}

	return diskIds, nil
}
