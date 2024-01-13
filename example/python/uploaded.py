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
import os

filename = "./example.txt"

# Ensure the file exists
if not os.path.exists(filename):
    raise Exception(f"File not found: {filename}")

# Read the file
with open(filename, 'rb') as file:
    file_bytes = file.read()

# URL to send the POST request
url = "https://beta.gpt4api.plus/standard/uploaded"

# Prepare the headers and data for the POST request
headers = {
    "Authorization": "Bearer <YOUR_ACCESS_TOKEN>",
    "Content-Type": "multipart/form-data"
}
data = {
    "conversation_id": "",
    "type": "my_files"
}
files = {
    "file": (os.path.basename(filename), file_bytes)
}

# Make the POST request
response = requests.post(url, headers=headers, data=data, files=files, timeout=8)

# Check for successful status code
if response.status_code != 200:
    raise Exception(f"Request failed with status: {response.status_code}, {response.reason}")

# Parse the response
result = response.json()
print("<Upload>:", result)
