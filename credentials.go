/*
Copyright 2023 The gpt4batch Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gpt4batch

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CredentialsResponse is the response o1f the credentials.
type CredentialsResponse struct {
	AccessToken string `json:"access_token"`
	ExpiredAt   int64  `json:"expired_at"`
	Username    string `json:"username"`
}

// ParseCredentials parses the credentials.
func ParseCredentials(path string) (string, error) {
	if govalidator.IsNull(path) {
		// Read the token. If the token is empty, panic.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		homeFilePath := filepath.Join(homeDir, ".gpt4api-sk.pub")
		token, err := ioutil.ReadFile(homeFilePath)
		if err != nil {
			return "", err
		}

		if string(token) == "" {
			return "", fmt.Errorf("%s has no token", homeFilePath)
		}

		var resp CredentialsResponse
		if err := json.Unmarshal(token, &resp); err != nil {
			return "", err
		}
		return resp.AccessToken, nil
	}

	// Read the token. If the token is empty, panic.
	token, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	if string(token) == "" {
		return "", fmt.Errorf("token is required")
	}

	var resp CredentialsResponse
	if err := json.Unmarshal(token, &resp); err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}
