package bakeryclient

import (
	"fmt"
	"net/http"
)

func (c *Client) DeleteImage(imageName string) error {
	fullUrl := c.url + "/bakeforms/" + imageName
	req, _ := http.NewRequest("DELETE", fullUrl, nil)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("DELETE bakeform returned status code %v", res.StatusCode)
	}

	return nil
}
