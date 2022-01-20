Start-Transcript "build.log" -Force

    # record path
    $path = Get-Location

    Set-Location $PSScriptRoot

        if (!(Test-Path "./build/")) {
            mkdir build
        }
        
        # config center service
        powershell -executionpolicy bypass -File ".\config_center\build_linux.ps1"

        Write-Output "> build config center service finished."

        # gateway service
        powershell -executionpolicy bypass -File ".\gateway\build_linux.ps1"

        Write-Output "> build gateway service finished."

        # user service
        powershell -executionpolicy bypass -File ".\user\build_linux.ps1"

        Write-Output "> build user service finished."

        # cloud file service
        powershell -executionpolicy bypass -File ".\cloud_file\build_linux.ps1"

        Write-Output "> build cloud file service finished."

        # thinking note service
        powershell -executionpolicy bypass -File ".\thinking_note\build_linux.ps1"

        Write-Output "> build thinking note service finished."

    # reset path
    Set-Location $path

    Write-Output "> build finished."

Stop-Transcript

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
