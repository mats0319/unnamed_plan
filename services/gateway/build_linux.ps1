# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../build/service_gateway") {
        Remove-Item "../build/service_gateway"
    }

    go build -o "service_gateway"

    Move-Item "service_gateway" -Destination "../build/service_gateway"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
