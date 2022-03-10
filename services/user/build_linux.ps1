# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../build/service_user/") {
        Remove-Item "../build/service_user/*"
    } else {
        mkdir "../build/service_user/"
    }

    go build -o "service_exec"

    Move-Item "service_exec" -Destination "../build/service_user/service_exec"

    Copy-Item "config_production.json" -Destination "../build/service_user/config.json"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
