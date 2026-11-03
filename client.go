package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client -
type Client struct {
	httpClient http.Client
}

// NewClient -
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}


func (c *Client) GetResponseById(ctx context.Context, baseURL *string, id string) (xmlBlob []byte, http_err error) {
	// Using "non-simple" HTTP GET method
	fullURL := *baseURL
	if len(id) > 0 {
		if strings.HasSuffix(fullURL, "/") {
			fullURL = fullURL + id
		} else {
			fullURL = fullURL + "/" + id
		}
	}
	// fmt.Printf("reading %s\n", fullURL)

    // create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "gator_go")

    // create a new client and make the request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		fmt.Printf("%s not found\n", fullURL)
		return []byte{}, nil
	} else if res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("HTTP failure (%d) for %s\n", res.StatusCode, fullURL)
	}

	xmlBlob, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return xmlBlob, nil
}

