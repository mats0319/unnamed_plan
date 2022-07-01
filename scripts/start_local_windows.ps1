$path = Get-Location

Set-Location $PSScriptRoot

    # core service
    Set-Location "../services/core/"

    go build -race -o service_exec.exe

    Start-Process ./service_exec.exe

    Set-Location $PSScriptRoot

    # gateway service
    Set-Location "../services/gateway/"

    go build -race -o service_exec.exe

    Start-Process ./service_exec.exe

    Set-Location $PSScriptRoot

    # business service 1: user service
    Set-Location "../services/1_user/"

    go build -race -o service_exec.exe

    Start-Process ./service_exec.exe

    Set-Location $PSScriptRoot

    # business servcie 2: cloud file service
    Set-Location "../services/2_cloud_file/"

    go build -race -o service_exec.exe

    Start-Process ./service_exec.exe

    Set-Location $PSScriptRoot

    # business service 3: note service
    Set-Location "../services/3_note/"

    go build -race -o service_exec.exe

    Start-Process ./service_exec.exe

    Set-Location $PSScriptRoot

    # business service 4: task service
    Set-Location "../services/4_task/"

    go build -race -o service_exec.exe

    Start-Process ./service_exec.exe

    Set-Location $PSScriptRoot

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
