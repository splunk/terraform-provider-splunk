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

default: build

fmt:
	go fmt ./...
	@terraform fmt -recursive

build:
	go build -o terraform-provider-splunk .

# Use external linker on macOS to avoid dyld "missing LC_UUID" abort (macOS 15+).
# TF_ACC= ensures acceptance tests are skipped (no API required).
test:
	TF_ACC= go test -ldflags="-linkmode=external" ./...

testacc:
	TF_ACC=1 go test -ldflags="-linkmode=external" ./... -v

init:
	@terraform init

plan:
	@terraform plan

apply:
	@terraform apply -auto-approve
