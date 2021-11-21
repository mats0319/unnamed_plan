# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../gateway_service") {
        Remove-Item "../gateway_service"
    }

    go build -o "gateway_service"

    Move-Item "gateway_service" -Destination "../gateway_service"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
