# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../admin_data") {
        Remove-Item "../admin_data"
    }

    go mod download

    go build -o "admin_data"

    Move-Item "admin_data" -Destination "../admin_data"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
