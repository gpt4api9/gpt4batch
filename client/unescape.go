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

package client

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// downloadReq is a request to download a file.
type downloadReq struct {
	Id     string
	Pid    string
	URL    string
	Prefix string
}

// DownloadUrlPath returns the path of the download URL.
func DownloadUrlPath(req *downloadReq) string {
	uf, err := url.PathUnescape(req.URL)
	if err != nil {
		return ""
	}

	uf = strings.ReplaceAll(uf, "&amp;", "&")
	uf = strings.ReplaceAll(uf, ";", "&")

	u, err := url.Parse(uf)
	if err != nil {
		return ""
	}

	queryParams, _ := url.ParseQuery(u.RawQuery)
	filename := queryParams.Get(" filename")
	if filename != "" {
		// TODO: check if filename is valid.
		list := make([]string, 0, 5)
		list = append(list, req.Prefix)
		list = append(list, req.Id)
		list = append(list, req.Pid)
		list = append(list, fmt.Sprintf("%d", time.Now().UnixNano()))
		list = append(list, filename)

		return strings.Join(list, "_")
	}
	return ""
}
