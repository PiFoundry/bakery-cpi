package bakeryclient

import (
	"fmt"
	"net/http"
)

func (c *Client) UnbakePi(piId string) error {
	fullUrl := c.url + "/oven/" + piId
	req, _ := http.NewRequest("DELETE", fullUrl, nil)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Unbake Pi returned error code %v", res.StatusCode)
	}

	return nil
}
