# set env in powershell, only apply in this cmd-line window
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"

$path = Get-Location

Set-Location $PSScriptRoot

    Set-Location "../services/task"

    if (Test-Path "../build/service_task/") {
        Remove-Item "../build/service_task/*"
    } else {
        mkdir "../build/service_task/"
    }

    go build -o "service_exec"

    Move-Item "service_exec" -Destination "../build/service_task/service_exec"

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
