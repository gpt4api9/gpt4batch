/*
Copyright 2024 The gpt4batch Authors.

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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gitlab.com/gpt4batch"
	"net/http"
	"time"
)

func main() {
	req := gpt4batch.OpenaiChatRequest{
		Message: "你是gpt3还是gpt-4",
		Model:   "gpt-4",
	}

	// 建议对话设置8分钟
	resp, err := resty.New().
		SetTimeout(8*time.Minute).
		R().
		EnableTrace().
		SetAuthToken("<YOUR_ACCESS_TOKEN>").
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post("https://beta.gpt4api.plus/standard/all-tools")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		panic(fmt.Sprintf(""))
	}

	var result gpt4batch.ChatResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		panic(err)
	}

	fmt.Println("<ChatResponse>: ", result)
}
