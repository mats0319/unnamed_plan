# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../build/service_cloud_file") {
        Remove-Item "../build/service_cloud_file"
    }

    go build -o "service_cloud_file"

    Move-Item "service_cloud_file" -Destination "../build/service_cloud_file"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
