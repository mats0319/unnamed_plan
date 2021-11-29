Start-Transcript "build.log" -Force

    # record path
    $path = Get-Location

    Set-Location $PSScriptRoot

        # gateway service
        powershell -executionpolicy bypass -File ".\services\gateway\build_linux.ps1"

        Write-Output "> build gateway service finished."

        # user service
        powershell -executionpolicy bypass -File ".\services\user\build_linux.ps1"

        Write-Output "> build user service finished."

        # cloud file service
        powershell -executionpolicy bypass -File ".\services\cloud_file\build_linux.ps1"

        Write-Output "> build cloud file service finished."

        # thinking note service
        powershell -executionpolicy bypass -File ".\services\thinking_note\build_linux.ps1"

        Write-Output "> build thinking note service finished."

    ###

    # admin_ui
#    Set-Location ".\admin_ui"
#    npm run link
#    npm run build-report
#
#    Set-Location ..
#
#    Write-Output "> build admin ui finished."

    # public_mobile
#    Set-Location ".\public_mobile"
#    npm run link
#    npm run build-report
#
#    Set-Location ..
#
#    Write-Output "> build public mobile finished."

    # public_ui
#    Set-Location ".\public_ui"
#    npm run link
#    npm run build-report
#
#    Set-Location ..
#
#    Write-Output "> build public ui finished."

    # reset path
    Set-Location $path

    Write-Output "> build finished."

Stop-Transcript

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
