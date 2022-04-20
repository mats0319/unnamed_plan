$path = Get-Location

Set-Location $PSScriptRoot

Set-Location "../services/cloud_file/rpc"

go test

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
