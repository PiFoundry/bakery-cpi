package bakeryclient

import (
	"bytes"
	"fmt"
	"net/http"
)

func (c *Client) UploadBytesAsFile(piId, filename string, filebytes []byte) error {
	fullUrl := fmt.Sprintf("%v/oven/%v/upload/%v", c.url, piId, filename)
	res, err := http.Post(fullUrl, "application/octet-stream", bytes.NewReader(filebytes))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("UploadBytesAsFile returned status code %v", res.StatusCode)
	}
	return nil
}
