package client

import (
	"bytes"
	"encoding/json"
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

// Netconf struct represents a Netconf Device Details
type Netconf struct {
	Name      string `json:"node-id"`
	IPAddress string `json:"netconf-node-topology:host"`
	Port      int    `json:"netconf-node-topology:port"`
	Username  string `json:"netconf-node-topology:username"`
	Password  string `json:"netconf-node-topology:password"`
}

// NetconfPayload struct
type NetconfPayload struct {
	Node []Netconf `json:"node"`
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

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil, err
	}

	log.Printf("[DEBUG] GET Body: ", string(bodyBytes))

	item := &NetconfPayload{}
	err = json.Unmarshal(bodyBytes, item)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil, err
	}

	log.Printf("[DEBUG] Parsed Body: ", item)

	return &item.Node[0], nil
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
func (c *Client) NetconfMount(n Netconf) error {

	payload := NetconfPayload{
		Node: []Netconf{n},
	}

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(payload)

	log.Printf("[DEBUG] Payload: ", buf)

	_, err = c.httpRequest(fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%s", n.Name), "PUT", buf)
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
