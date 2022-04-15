# build: build linux exec of back-end services
# start: local start back-end services
.PHONY: build start

build: scripts/build_linux.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux.ps1"

start:
	make start_local_gateway
	make start_local_user
	make start_local_cloud_file
	make start_local_note
	make start_local_task

# build one service
.PHONY: build_linux_gateway build_linux_user build_linux_cloud_file build_linux_note build_linux_task

build_linux_gateway: scripts/build_linux_gateway.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux_gateway.ps1"

build_linux_user: scripts/build_linux_user.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux_user.ps1"

build_linux_cloud_file: scripts/build_linux_cloud_file.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux_cloud_file.ps1"

build_linux_note: scripts/build_linux_note.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux_note.ps1"

build_linux_task: scripts/build_linux_task.ps1
	powershell -executionpolicy bypass -File "./scripts/build_linux_task.ps1"

# start one service
.PHONY: start_local_gateway start_local_user start_local_cloud_file start_local_note start_local_task

start_local_gateway: scripts/start_local_gateway.ps1
	powershell -executionpolicy bypass -File "./scripts/start_local_gateway.ps1"

start_local_user: scripts/start_local_user.ps1
	powershell -executionpolicy bypass -File "./scripts/start_local_user.ps1"

start_local_cloud_file: scripts/start_local_cloud_file.ps1
	powershell -executionpolicy bypass -File "./scripts/start_local_cloud_file.ps1"

start_local_note: scripts/start_local_note.ps1
	powershell -executionpolicy bypass -File "./scripts/start_local_note.ps1"

start_local_task: scripts/start_local_task.ps1
	powershell -executionpolicy bypass -File "./scripts/start_local_task.ps1"
