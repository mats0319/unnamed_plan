.phony: build_linux_services

#<target> : <prerequisites>
#[tab]  <commands>

build_linux_services: services/scripts/build_linux_services.ps1
	powershell -executionpolicy bypass -File "./services/scripts/build_all_linux.ps1"
