package gatewayapi

import (
	"errors"
	"strconv"
	"strings"
)

type Usage struct {
	Countries map[string]int `json:"countries"`
	Currency  string         `json:"currency"`
	TotalCost float64        `json:"total_cost"`
}

type MtSmsResponse struct {
	IDs   []int `json:"ids"`
	Usage Usage `json:"usage"`
}

type ErrorResponse struct {
	Code         string   `json:"code"`
	IncidentUuid string   `json:"incident_uuid"`
	Message      string   `json:"message"`
	Variables    []string `json:"variables"`
}

func (r ErrorResponse) Error() error {
	msg := r.Message
	if r.Code != "" {
		msg = "code: " + r.Code + "; " + msg
	}
	for i := range r.Variables {
		msg = strings.Replace(msg, "%"+strconv.FormatInt(int64(i), 10), r.Variables[i], 1)
	}
	return errors.New(msg)
}
