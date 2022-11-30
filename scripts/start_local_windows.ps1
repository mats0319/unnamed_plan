$path = Get-Location

Set-Location $PSScriptRoot

    # core service
    Set-Location "../services/core/"

    go build -race -o service_core.exe

    Start-Process ./service_core.exe

    Set-Location $PSScriptRoot

    # gateway service
    Set-Location "../services/gateway/"

    go build -race -o service_gateway.exe

    Start-Process ./service_gateway.exe

    Set-Location $PSScriptRoot

    # business service 1: user service
    Set-Location "../services/1_user/"

    go build -race -o service_1_user.exe

    Start-Process ./service_1_user.exe

    Set-Location $PSScriptRoot

    # business servcie 2: cloud file service
    Set-Location "../services/2_cloud_file/"

    go build -race -o service_2_cloud_file.exe

    Start-Process ./service_2_cloud_file.exe

    Set-Location $PSScriptRoot

    # business service 3: note service
    Set-Location "../services/3_note/"

    go build -race -o service_3_note.exe

    Start-Process ./service_3_note.exe

    Set-Location $PSScriptRoot

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
