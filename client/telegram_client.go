package client

import (
	"bytes"
	"encoding/json"
	"errors"
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
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.New("unable to create a new request")
	}

	response, err := r.clients.Do(request)
	if err != nil {
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

func FetchAdventurerInfo(ClientID string) (*domain.Response, error) {
	URL := "https://cryptoguild.herokuapp.com/AdventurerStatus/" + ClientID
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cResp domain.Response
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return nil, err
	}
	//Invoke the text output function & return it with nil as the error value
	return &cResp, nil
}

func getClient(url string, cBody interface{}, cResp interface{}) (interface{}, error) {

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(cBody)
	if err != nil {
		return nil, err
	}

	// Send a post request with your token
	// "https://api.telegram.org/bot1913861473:AAGT0ranx9RBMrtRVzrLx5PYiakOsNH6VOE/sendMessage"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status" + resp.Status)
	}

	defer resp.Body.Close()

	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return nil, err
	}

	return &cResp, nil
}
