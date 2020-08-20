package awair

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	token string

	Options struct {
		PreferFahrenheit bool
	}
}

func NewClientFromBearerToken(token string) *Client {
	return &Client{token: token}
}

func (c *Client) do(ctx context.Context, method, url string, formdata map[string]string, dest interface{}) error {
	var body io.Reader
	if formdata != nil {
		form, err := json.Marshal(formdata)
		if err != nil {
			return err
		}
		body = bytes.NewReader(form)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if dest == nil {
		return nil
	}

	return json.Unmarshal(data, dest)
}

func (c *Client) get(ctx context.Context, url string, dest interface{}) error {
	return c.do(ctx, "GET", url, nil, dest)
}

func (c *Client) put(ctx context.Context, url string, formdata map[string]string) error {
	return c.do(ctx, "PUT", url, formdata, nil)
}

func (c *Client) GetDevices(ctx context.Context) ([]*Device, error) {
	var wrapper struct {
		Devices []*Device `json:"devices"`
	}
	err := c.get(ctx, "https://developer-apis.awair.is/v1/users/self/devices", &wrapper)
	if err != nil {
		return nil, err
	}

	for _, d := range wrapper.Devices {
		d.c = c
	}

	return wrapper.Devices, nil
}
