# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../build/service_cloud_file/") {
        Remove-Item "../build/service_cloud_file/*"
    } else {
        mkdir "../build/service_cloud_file/"
    }

    go build -o "service_exec"

    Move-Item "service_exec" -Destination "../build/service_cloud_file/service_exec"

    Copy-Item "config_production.json" -Destination "../build/service_cloud_file/config.json"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
