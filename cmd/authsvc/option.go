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

package authsvc

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

// Option represents the option of the token service.
type Option struct {
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
	TTL      int    `json:"ttl"`
}

// Validate validates the option.
func (o Option) Validate() error {
	if govalidator.IsNull(o.URL) {
		return errors.New("url is required")
	}

	if !govalidator.IsEmail(o.Email) {
		return errors.New("email is required")
	}

	if govalidator.IsNull(o.Password) {
		return errors.New("password is required")
	}
	return nil
}
