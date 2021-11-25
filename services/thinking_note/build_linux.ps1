# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../build/service_thinking_note") {
        Remove-Item "../build/service_thinking_note"
    }

    go build -o "service_thinking_note"

    Move-Item "service_thinking_note" -Destination "../build/service_thinking_note"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
