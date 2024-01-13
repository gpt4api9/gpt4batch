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

url = "https://beta.gpt4api.plus/standard/chat/uploaded"

payload = {'type': 'my_files',
'conversation_id': ''}
files=[
  ('file',('1.txt',open('1.txt','rb'),'text/plain'))
]
headers = {
  'Authorization': 'Bearer <YOUR-ACCESS_TOKEN>'
}

response = requests.request("POST", url, headers=headers, data=payload, files=files)

print(response.text)
