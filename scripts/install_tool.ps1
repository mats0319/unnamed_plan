$path = Get-Location

Set-Location $PSScriptRoot

    Set-Location "../services/tools/protoc_ts"

    go install

Set-Location $path