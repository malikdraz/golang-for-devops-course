package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, loginURL, password string) (string, error) {
	LoginRequest := LoginRequest{
		Password: password,
	}
	body, err := json.Marshal(LoginRequest)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", err)
	}
	response, err := client.Post(loginURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", fmt.Errorf("Post error: %s", err)
	}

	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Invalid output (HTTP Code %d): %s\n", response.StatusCode, string(resBody))
	}

	var LoginResponse LoginResponse

	if !json.Valid(resBody) {
		return "", RequestError{
			Err:      fmt.Sprintf("Response is not a json"),
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
		}
	}

	err = json.Unmarshal(resBody, &LoginResponse)
	if err != nil {
		return "", RequestError{
			Err:      fmt.Sprintf("Page unmarshal error: %s", err),
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
		}
	}

	//if LoginResponse.Token == "" {
	//	return "", RequestError{
	//		HTTPCode: response.StatusCode,
	//		Body:     string(resBody),
	//		Err:      "Empty token replied",
	//	}
	//}

	return LoginResponse.Token, nil
}
