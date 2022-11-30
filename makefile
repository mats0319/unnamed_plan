# process_code: use 'go fmt' and 'go vet' process go code
# build: build linux exec of back-end services
# start: local start back-end services，with '-race' flag
# install: install protoc-ts tool
# protoc: generate go and ts code from proto file
.PHONY: process_code build start install protoc

process_code: scripts/process_code.ps1
	powershell -executionpolicy bypass -File "./scripts/process_code.ps1"

build: scripts/build_linux.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux.ps1"

start: scripts/start_local_windows.ps1
	powershell -executionpolicy bypass -File "./scripts/start_local_windows.ps1"

install: scripts/install_tool.ps1
	powershell -executionpolicy bypass -File "./scripts/install_tool.ps1"

protoc: services/shared/proto/generate_code.ps1
	powershell -executionpolicy bypass -File "./services/shared/proto/generate_code.ps1"
