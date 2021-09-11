# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

if (Test-Path "../public_data") {
    Remove-Item "../public_data"
}

go mod download

go build -o "public_data"

Move-Item "public_data" -Destination "../public_data"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
