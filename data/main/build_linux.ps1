# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

Set-Location $PSScriptRoot
go build

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
