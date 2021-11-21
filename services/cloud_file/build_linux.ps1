# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../cloud_file_service") {
        Remove-Item "../cloud_file_service"
    }

    go build -o "cloud_file_service"

    Move-Item "cloud_file_service" -Destination "../cloud_file_service"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
