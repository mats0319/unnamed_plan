# process_code: use 'go fmt' and 'go vet' process go code
# build: build linux exec of back-end services
# start: local start back-end services，with '-race' flag
.PHONY: process_code build start

process_code: scripts/process_code.ps1
	powershell -executionpolicy bypass -File "./scripts/process_code.ps1"

build: scripts/build_linux.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux.ps1"

start: scripts/start_local_windows.ps1
	powershell -executionpolicy bypass -File "./scripts/start_local_windows.ps1"
