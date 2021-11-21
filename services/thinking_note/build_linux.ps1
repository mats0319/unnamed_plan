# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    if (Test-Path "../thinking_note_service") {
        Remove-Item "../thinking_note_service"
    }

    go build -o "thinking_note_service"

    Move-Item "thinking_note_service" -Destination "../thinking_note_service"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
