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
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/gpt4batch"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"gitlab.com/gpt4batch/log"
)

// credentials is the user credentials.
type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TTL      int    `json:"ttl"`
}

// NewAuthenticationCommand creates a new auth token command.
func NewAuthenticationCommand(ctx context.Context) *cobra.Command {
	// create a new option instance. this is passed to the command handler.
	var (
		option Option
		logger = log.New(log.InfoLevel)
	)

	rootCmd := &cobra.Command{
		Use:  "authsvc",
		Args: cobra.NoArgs,
		Short: "Retrieve the GPT4Batch token for authentication. " +
			"and only one valid token will be retained. " +
			"please use this command with caution.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// validate the option.
			if err := option.Validate(); err != nil {
				return err
			}

			logger.
				WithField("email", option.Email).
				WithField("password", option.Password).
				Info("config")

			// set the default value of the option.
			// store token in the home directory.
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			// homeDir is the user's home directory. + ".gpt4api-sk.pub"
			homeFilePath := filepath.Join(homeDir, ".gpt4api-sk.pub")

			logger.
				WithField("dir", homeDir).
				WithField("pub", homeFilePath).
				Info("homeDir")

			// send the request to the server.
			// Use the default if available
			resp, err := resty.New().R().
				SetHeader("Content-Type", "application/json").
				SetBody(credentials{
					Email:    option.Email,
					Password: option.Password,
					TTL:      option.TTL,
				}).
				Post(option.URL)
			if err != nil {
				return err
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("%s", resp.Status())
			}

			var oct gpt4batch.CredentialsResponse
			if err := json.Unmarshal(resp.Body(), &oct); err != nil {
				return err
			}

			// marshal the token.
			result, _ := json.Marshal(oct)

			// write the token to the home directory.
			// todo: use the viper to store the token.
			if err = os.WriteFile(homeFilePath, result, 0644); err != nil {
				return err
			}

			logger.
				WithField("expired_at", oct.ExpiredAt).
				WithField("username", oct.Username).
				WithField("token", oct.AccessToken).
				Info("ok")
			return nil
		},
	}

	rootCmd.Flags().StringVarP(&option.URL, "url", "u", "https://beta.gpt4api.shop/console/access_token", "设置获取Ak调用服务地址.")
	rootCmd.Flags().StringVarP(&option.Email, "email", "e", "", "输入账号邮箱.如果没注册可访问官网.https://gpt4api.shop.")
	rootCmd.Flags().StringVarP(&option.Password, "password", "p", "", "输入账号密码,如果没注册可访问官网.https://gpt4api.shop.")
	rootCmd.Flags().IntVarP(&option.TTL, "ttl", "t", 24*60*60, "设置AccessToken过期时间，默认是60天.")
	return rootCmd
}
