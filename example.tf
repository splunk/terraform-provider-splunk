# Copyright 2015 Container Solutions
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

provider "splunk" {
  url                  = "localhost:8089"
  username             = "admin"
  password             = "changeme"
  insecure_skip_verify = true
}

resource "splunk_global_http_event_collector" "http" {
  disabled    = false
  enable_ssl  = true
  port        = 8088
}

resource "splunk_input_http_event_collector" "hec-token" {
  name       = "new-token"
  index      = "main"
  source     = "new-source"
  sourcetype = "new-sourcetype"
  disabled   = false
  use_ack    = false

  depends_on = ["splunk_global_http_event_collector.http"]
}
