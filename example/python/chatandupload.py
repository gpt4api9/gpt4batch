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

# File upload section
filename = "./example.txt"

# Ensure the file exists
if not os.path.exists(filename):
    raise Exception(f"File not found: {filename}")

# Read the file
with open(filename, 'rb') as file:
    file_bytes = file.read()

# URL for the file upload
upload_url = "https://beta.gpt4api.plus/standard/uploaded"

# Prepare headers and data for the file upload
upload_headers = {
    "Authorization": "Bearer <YOUR_ACCESS_TOKEN>",
    "Content-Type": "multipart/form-data"
}
upload_data = {
    "conversation_id": "",
    "type": "my_files"
}
files = {
    "file": (os.path.basename(filename), file_bytes)
}

# Make the file upload POST request
upload_response = requests.post(upload_url, headers=upload_headers, data=upload_data, files=files, timeout=8)

# Check the response and parse it
if upload_response.status_code != 200:
    raise Exception(f"File upload failed: {upload_response.status_code}, {upload_response.reason}")
upload_result = upload_response.json()

# Chat request section
chat_req_data = {
    "gizmo_id": "g-HMNcP6w7d",
    "message": "你是gpt3还是gpt-4",
    "model": "gpt-4",
    "attachments": [upload_result['attachment']],
    "parts": [upload_result['part']]
}

# URL for the chat request
chat_url = "https://beta.gpt4api.plus/standard/all-tools"

# Headers for the chat request
chat_headers = {
    "Authorization": "Bearer <YOUR_ACCESS_TOKEN>",
    "Content-Type": "application/json"
}

# Make the chat POST request
chat_response = requests.post(chat_url, headers=chat_headers, json=chat_req_data, timeout=8*60)

# Check the response and parse it
if chat_response.status_code != 200:
    raise Exception(f"Chat request failed: {chat_response.status_code}, {chat_response.reason}")
chat_result = chat_response.json()

# Print the chat response
print("<ChatResponse>:", chat_result)
