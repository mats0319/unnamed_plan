# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (!(Test-Path "../build/service_gateway/")) {
        mkdir "../build/service_gateway/"
    }

    if (Test-Path "../build/service_gateway/service_exec") {
        Remove-Item "../build/service_gateway/service_exec"
    }

    go build -o "service_exec"

    Move-Item "service_exec" -Destination "../build/service_gateway/service_exec"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
