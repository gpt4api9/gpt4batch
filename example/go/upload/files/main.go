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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gitlab.com/gpt4batch"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	filename := "./example.txt"

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	resp, err := resty.New().
		SetTimeout(8*time.Second).
		R().
		EnableTrace().
		SetAuthToken("<YOUR_ACCESS_TOKEN>").
		SetHeader("Content-Type", "multipart/form-data").
		SetFormData(map[string]string{
			"conversation_id": "",
			"type":            gpt4batch.MyFiles,
		}).
		SetFileReader("file", filepath.Base(filename), bytes.NewReader(fileBytes)).
		Post("https://beta.gpt4api.plus/standard/uploaded")

	if err != nil {
		panic(err)
	}

	if resp.StatusCode() != http.StatusOK {
		panic(resp.Status())
	}

	var result gpt4batch.UploadResponse
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		panic(err)
	}

	fmt.Println("<Upload>: ", result)
}
