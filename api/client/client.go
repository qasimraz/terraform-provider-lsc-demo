package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Client holds all of the information required to connect to a controller
type Client struct {
	hostname   string
	port       int
	authToken  string
	httpClient *http.Client
}

// NewClient returns a new Lumina SDN controller client
func NewClient(hostname string, port int, token string) *Client {
	return &Client{
		hostname:   hostname,
		port:       port,
		authToken:  token,
		httpClient: &http.Client{},
	}
}

// GetNetconf gets a generic netconf endpoint with a url from the controller
func (c *Client) GetNetconf(url string) ([]byte, error) {
	body, err := c.httpRequest(url, "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil, err
	}

	log.Printf("[DEBUG] GET Body: ", string(bodyBytes))

	return bodyBytes, nil
}

// DeleteNetconf deletes a generic netconf endpoint from the controller
func (c *Client) DeleteNetconf(url string) error {
	_, err := c.httpRequest(url, "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}

// PutNetconf puts a netconf payload at a specific url mount point
func (c *Client) PutNetconf(url string, payloadBody bytes.Buffer) error {
	_, err := c.httpRequest(url, "PUT", payloadBody)
	if err != nil {
		return err
	}
	return nil
}

// httpRequest calls generic HTTP requests
func (c *Client) httpRequest(path string, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authToken)
	switch method {
	case "GET":
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}

	log.Printf("[DEBUG] API call: ", req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%s:%v/%s", c.hostname, c.port, path)
}
