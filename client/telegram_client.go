package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gunnsteinn/telegramBot/domain"
	"io/ioutil"
	"net/http"
)

var (
	ResponseClient responseClientInterface = &responseClient{}
)

type responseClient struct {
	clients http.Client
}

type responseClientInterface interface {
	Get(url string) (*domain.Response, error)
	Post(url string, body interface{}) (*domain.Response, error)
}

func (r *responseClient) Get(url string) (*domain.Response, error) {
	return r.do(http.MethodGet, url, nil)
}

func (r *responseClient) Post(url string, body interface{}) (*domain.Response, error) {
	return r.do(http.MethodPost, url, body)
}

func (r *responseClient) do(method string, url string, body interface{}) (*domain.Response, error) {

	requestBody, err := getRequestBody(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.New("unable to create a new request")
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := r.clients.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusInternalServerError {
		defer response.Body.Close()
		errBody, _ := responseBodyAsString(response)
		fmt.Println("internal server error. " + errBody)
		return nil, err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	finalResponse := domain.Response{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Headers:    response.Header,
		Body:       responseBody,
	}
	return &finalResponse, nil
}

func getRequestBody(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	return json.Marshal(body)
}

func responseBodyAsString(resp *http.Response) (string, error) {
	if resp.Body == nil {
		return "", nil
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	return string(bytes), err
}
