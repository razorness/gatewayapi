package gatewayapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/mrjones/oauth"
)

const (
	MtSmsEndpoint = "https://gatewayapi.com/rest/mtsms"
)

var (
	client *Client
	once   sync.Once
)

type Client struct {
	key    string
	secret string
	client *http.Client
}

func NewClient(key, secret string) *Client {

	once.Do(func() {
		consumer := oauth.NewConsumer(key, secret, oauth.ServiceProvider{})
		c, err := consumer.MakeHttpClient(&oauth.AccessToken{})
		if err != nil {
			panic(err)
		}
		c.Timeout = 30 * time.Second

		client = &Client{
			key:    key,
			secret: secret,
			client: c,
		}
	})

	return client
}

func (c *Client) SendSms(sms *SMS) (*MtSmsResponse, error) {

	resp, err := c.Do(http.MethodPost, MtSmsEndpoint, sms)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {

		var res MtSmsResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil

	}

	var er ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, err
	}

	return nil, er.Error()

}

func (c *Client) Do(method, url string, body interface{}) (*http.Response, error) {

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "razorness/gatewayapi")
	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)

}
