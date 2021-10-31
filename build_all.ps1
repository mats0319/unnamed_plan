Start-Transcript "build.log" -Force

    # record path
    $path = Get-Location

    Set-Location $PSScriptRoot

    # admin_data
    powershell -executionpolicy bypass -File ".\admin_data\main\build_linux.ps1"

    Write-Output "> build admin data finished."

    # public_data
    powershell -executionpolicy bypass -File ".\public_data/main\build_linux.ps1"

    Write-Output "> build public data finished."

    ###

    # admin_ui
    Set-Location ".\admin_ui"
    npm run link
    npm run build-report

    Set-Location ..

    Write-Output "> build admin ui finished."

    # public_mobile
    Set-Location ".\public_mobile"
    npm run link
    npm run build-report

    Set-Location ..

    Write-Output "> build public mobile finished."

    # public_ui
    Set-Location ".\public_ui"
    npm run link
    npm run build-report

    Set-Location ..

    Write-Output "> build public ui finished."

    # reset path
    Set-Location $path

    Write-Output "> build finished."

Stop-Transcript
