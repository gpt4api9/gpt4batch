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
	"bufio"
	"encoding/json"
	"fmt"
	"gitlab.com/gpt4batch"
	"os"
)

func main() {
	for i := 0; i < 1000; i++ {
		in := &gpt4batch.In{
			ID: fmt.Sprintf("id: %d", i),
		}

		for j := 0; j < 1; j++ {
			in.Asks = append(in.Asks, &gpt4batch.Ask{
				ID:      fmt.Sprintf("pid: %d", j),
				Content: "你是gpt3还是gpt-4",
			})
		}

		Append("example.jsonl", in)
	}
}

func Append(outFilename string, item interface{}) error {
	file, err := os.OpenFile(outFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonStr, err := json.Marshal(item)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(file)
	w.Write(jsonStr)
	w.WriteString("\n")
	return w.Flush()
}
