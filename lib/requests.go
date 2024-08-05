package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Target struct {
	Url string
}

func (target *Target) SessionInit() *http.Client {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	target.Url = strings.TrimSuffix(target.Url, "/")
	target.Url = fmt.Sprintf("%s/robots.txt", target.Url)
	return client
}

func (target *Target) SendGET(session http.Client) (*http.Response, error) {
	response, err := session.Get(target.Url)
	if err != nil {
		return nil, errors.New("GET request failed: " + err.Error())
	}
	if response.StatusCode != http.StatusOK {
		return response, errors.New("received non-200 status code: " + http.StatusText(response.StatusCode))
	}
	return response, nil
}

func (target *Target) HttpStatusCode(session *http.Client) int {
	response, err := target.SendGET(*session)
	if err != nil {
		fmt.Println(err)
	}
	return response.StatusCode
}

func (target *Target) GetRobotsFile(session *http.Client) (string, error) {
	response, err := target.SendGET(*session)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("failed to read response body: " + err.Error())
	}
	defer response.Body.Close()
	return string(body), nil
}
