Start-Transcript "test.log" -Force

    # record path
    $path = Get-Location

    Set-Location $PSScriptRoot

        # test user service
        powershell -executionpolicy bypass -File "./test_user.ps1"

        # test cloud file service
        powershell -executionpolicy bypass -File "./test_cloud_file.ps1"

        # test note service
        powershell -executionpolicy bypass -File "./test_note.ps1"

        # test task service
        powershell -executionpolicy bypass -File "./test_task.ps1"

    # reset path
    Set-Location $path

Stop-Transcript

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
