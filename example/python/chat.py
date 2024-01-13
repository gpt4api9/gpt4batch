# Copyright 2024 The gpt4batch Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


import requests
import json

# Define the request data
req_data = {
    "message": "你是gpt3还是gpt-4",
    "model": "gpt-4"
}

# URL for the POST request
url = "https://beta.gpt4api.plus/standard/all-tools"

# Set up headers
headers = {
    "Authorization": "Bearer <YOUR_ACCESS_TOKEN>",
    "Content-Type": "application/json"
}

# Make the POST request with a timeout of 8 minutes
response = requests.post(url, headers=headers, json=req_data, timeout=8*60)

# Check for successful status code
if response.status_code != 200:
    raise Exception(f"Request failed with status: {response.status_code}, {response.reason}")

# Parse the JSON response
result = response.json()
print("<ChatResponse>:", result)
