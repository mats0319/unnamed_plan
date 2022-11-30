Start-Transcript "build_linux_services.log" -Force

    # record path
    $path = Get-Location

    Set-Location $PSScriptRoot

    # prepare
    Set-Location "../services"

    go mod tidy

    if (Test-Path "./build/") {
        Remove-Item "./build/*" -Recurse
    } else {
        mkdir "./build"
    }

    mkdir "./build/service_core/"
    mkdir "./build/service_gateway/"
    mkdir "./build/service_1_user/"
    mkdir "./build/service_2_cloud_file/"
    mkdir "./build/service_3_note/"
    mkdir "./build/service_4_task/"

    Copy-Item "config.json" -Destination "./build/config.json"

    Set-Location $PSScriptRoot

    Copy-Item "upgrade_services.sh" -Destination "../services/build/upgrade_services.sh"

    # set env in powershell, only apply in this cmd-line window
    $env:CGO_ENABLED=0
    $env:GOOS="linux"
    $env:GOARCH="amd64"

        # services

        # core service
        Set-Location "../services/core/"

        go build -o "service_exec"

        Move-Item "service_exec" -Destination "../build/service_core/unnamed_plan_service_core"
        Copy-Item "config_production.json" -Destination "../build/service_core/config.json"

        Set-Location $PSScriptRoot

        Write-Output "> build core service finished"

        # gateway service
        Set-Location "../services/gateway/"

        go build -o "service_exec"

        Move-Item "service_exec" -Destination "../build/service_gateway/unnamed_plan_service_gateway"

        Set-Location $PSScriptRoot

        Write-Output "> build gateway service finished"

        # business service 1: user service
        Set-Location "../services/1_user/"

        go build -o "service_exec"

        Move-Item "service_exec" -Destination "../build/service_1_user/unnamed_plan_service_1_user"

        Set-Location $PSScriptRoot

        Write-Output "> build user service finished"

        # business service 2: cloud file service
        Set-Location "../services/2_cloud_file/"

        go build -o "service_exec"

        Move-Item "service_exec" -Destination "../build/service_2_cloud_file/unnamed_plan_service_2_cloud_file"

        Set-Location $PSScriptRoot

        Write-Output "> build cloud file service finished"

        # business service 3: note service
        Set-Location "../services/3_note/"

        go build -o "service_exec"

        Move-Item "service_exec" -Destination "../build/service_3_note/unnamed_plan_service_3_note"

        Set-Location $PSScriptRoot

        Write-Output "> build note service finished"

    # reset path
    Set-Location $path

Stop-Transcript

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
