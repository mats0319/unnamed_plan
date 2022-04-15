Start-Transcript "build_linux_services.log" -Force

    # record path
    $path = Get-Location

    Set-Location $PSScriptRoot

    Set-Location "../services"

    go mod tidy

    if (!(Test-Path "./build/")) {
        mkdir "./build"
    }

    Set-Location $PSScriptRoot

        # gateway service
        powershell -executionpolicy bypass -File "./build_linux_gateway.ps1"

        Write-Output "> build gateway service finished."

        # user service
        powershell -executionpolicy bypass -File "./build_linux_user.ps1"

        Write-Output "> build user service finished."

        # cloud file service
        powershell -executionpolicy bypass -File "./build_linux_cloud_file.ps1"

        Write-Output "> build cloud file service finished."

        # thinking note service
        powershell -executionpolicy bypass -File "./build_linux_note.ps1"

        Write-Output "> build thinking note service finished."

        # task service
        powershell -executionpolicy bypass -File "./build_linux_task.ps1"

        Write-Output "> build task service finished."

    # reset path
    Set-Location $path

    Write-Output "> build finished."

Stop-Transcript

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
