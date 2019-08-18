package yeelight

import (
	"errors"
	"fmt"
)

type Notification struct {
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

type Response interface {
	id() int
	ok() error
}

type OKResponse struct {
	ID     int      `json:"id"`
	Result []string `json:"result"`
}

func (r *OKResponse) id() int {
	return r.ID
}

func (r *OKResponse) ok() error {
	return nil
}

type ERRResponse struct {
	ID    int                    `json:"id"`
	Error map[string]interface{} `json:"error"`
}

func (r *ERRResponse) id() int {
	return r.ID
}

func (r *ERRResponse) ok() error {
	return errors.New(fmt.Sprintf("%v", r.Error))
}
