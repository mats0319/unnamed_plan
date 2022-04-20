$path = Get-Location

Set-Location $PSScriptRoot

Set-Location "../services/note/rpc"

go test

Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
