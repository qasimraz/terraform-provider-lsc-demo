package client

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// Client holds all of the information required to connect to a controller
type Client struct {
	hostname   string
	port       int
	authToken  string
	httpClient *http.Client
}

// Netconf struct represents a Netconf Device Details
type Netconf struct {
	IPAddress string `json:"host"`
	Name      string `json:"name"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
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

// GetNetconfMount gets a netconf mount with a specific name from the controller
func (c *Client) GetNetconfMount(name string) (*Netconf, error) {
	body, err := c.httpRequest(fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%v", name), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	item := &Netconf{}
	err = json.NewDecoder(body).Decode(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteNetconfMount removes a netconf mount point from the controller
func (c *Client) DeleteNetconfMount(name string) error {
	_, err := c.httpRequest(fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%v", name), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}

// NetconfMount creates a new device mount point
func (c *Client) NetconfMount(n *Netconf) error {
	payload := (`{
		"node": [
			{
				"node-id": "cisco1",
				"netconf-node-topology:host": "207.226.253.52",
				"netconf-node-topology:password": "root",
				"netconf-node-topology:username": "root",
				"netconf-node-topology:port": 830,
				"netconf-node-topology:tcp-only": false,
				"netconf-node-topology:keepalive-delay": 60
			}
		]
	}`)

	buf := bytes.Buffer{}
	err := xml.NewEncoder(&buf).Encode(payload)

	_, err = c.httpRequest(fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%v", n.Name), "PUT", buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) httpRequest(path, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authToken)
	switch method {
	case "GET":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
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
