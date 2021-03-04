package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

type apiTest struct {
	baseURL  string
	response *http.Response
	respErr  error
	client   *http.Client
}

func (a *apiTest) iRequestRESTEndpointWithMethodAndPathWithBody(method, path string, payload *messages.PickleStepArgument_PickleDocString) error {
	fullURL := a.baseURL + path
	req, err := http.NewRequest(method, fullURL, bytes.NewBufferString(payload.Content))
	if err != nil {
		return err
	}
	req.Body.Close()

	a.response, a.respErr = a.client.Do(req)
	if a.respErr != nil {
		return a.respErr
	}
	defer a.response.Body.Close()

	return nil
}

func (a *apiTest) thereShouldBeNoError() error {
	if a.respErr != nil {
		return a.respErr
	}
	return nil
}

func (a *apiTest) theResponseStatusCodeShouldBe(statusCode int) error {
	if a.response.StatusCode != statusCode {
		return fmt.Errorf("expected the status code '%v' in response to have value , but actually found this: '%v'", statusCode, a.response.StatusCode)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	test := apiTest{baseURL: "http://localhost:3000", client: &http.Client{}}

	s.Step(`^I request REST endpoint with method "([^"]*)" and path "([^"]*)" with body$`, test.iRequestRESTEndpointWithMethodAndPathWithBody)
	s.Step(`^the response status code should be "([^"]*)"$`, test.theResponseStatusCodeShouldBe)
	s.Step(`^there should be no error$`, test.thereShouldBeNoError)
}
