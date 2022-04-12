Start-Transcript "log.log" -Force

    # record path
    $path = Get-Location

    Set-Location $PSScriptRoot

        go mod tidy

        if (!(Test-Path "../build/")) {
            mkdir build
        }

        # gateway service
        powershell -executionpolicy bypass -File "..\gateway\build_linux.ps1"

        Write-Output "> build gateway service finished."

        # user service
        powershell -executionpolicy bypass -File "..\user\build_linux.ps1"

        Write-Output "> build user service finished."

        # cloud file service
        powershell -executionpolicy bypass -File "..\cloud_file\build_linux.ps1"

        Write-Output "> build cloud file service finished."

        # thinking note service
        powershell -executionpolicy bypass -File "..\note\build_linux.ps1"

        Write-Output "> build thinking note service finished."

        # task service
        powershell -executionpolicy bypass -File "..\task\build_linux.ps1"

        Write-Output "> build task service finished."

    # reset path
    Set-Location $path

    Write-Output "> build finished."

Stop-Transcript

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
