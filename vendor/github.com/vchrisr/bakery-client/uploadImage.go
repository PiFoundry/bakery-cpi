package bakeryclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type uploadResponse struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func (c *Client) UploadImage(imagePath, imageName string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fullUrl := c.url + "/bakeforms/" + imageName
	res, err := http.Post(fullUrl, "application/x-raw-disk-image", file)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("POST bakeform returned status code %v", res.StatusCode)
	}

	var parsedResponse uploadResponse
	message, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(message, &parsedResponse)
	if err != nil {
		return "", err
	}

	return parsedResponse.Name, nil
}
