package gatewayapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/mrjones/oauth"
)

const (
	MtSmsEndpoint = "https://gatewayapi.com/rest/mtsms"
)

var (
	client                   *Client
	once                     sync.Once
	UnauthorizedError        = errors.New("ie. invalid API key or signature")
	ForbiddenError           = errors.New("ie. unauthorized ip address")
	UnprocessableEntityError = errors.New("invalid json request body")
)

type Client struct {
	key     string
	secret  string
	client  *http.Client
	IsDebug bool
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
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, UnauthorizedError
	case http.StatusUnprocessableEntity:
		return nil, UnprocessableEntityError
	case http.StatusForbidden:
		return nil, ForbiddenError
	case http.StatusOK:

		var res MtSmsResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil

	default:

		var er ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
			return nil, err
		}

		return nil, er.Error()

	}

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

	if c.IsDebug {
		debuf, bodyErr := ioutil.ReadAll(req.Body)
		if bodyErr != nil {
			log.Println("gatewayapi bodyErr", bodyErr.Error())
		}

		rdr1 := ioutil.NopCloser(bytes.NewBuffer(debuf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(debuf))
		log.Println("gatewayapi request body:")
		log.Println(rdr1)
		req.Body = rdr2
	}

	return c.client.Do(req)

}
