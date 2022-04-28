# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

Set-Location "../services/config_center"

    if (Test-Path "../build/service_config_center/") {
        Remove-Item "../build/service_config_center/*"
    } else {
        mkdir "../build/service_config_center/"
    }

    go build -o "service_exec"

    Move-Item "service_exec" -Destination "../build/service_config_center/service_exec"

    Copy-Item "config_production.json" -Destination "../build/service_config_center/config.json"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
