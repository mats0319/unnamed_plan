$path = Get-Location

Set-Location $PSScriptRoot

Set-Location "../services/"

    # go fmt
    Write-Output "> go fmt start"

    gofmt -w -l -s .

    Write-Output "> go fmt finished"

    # go vet
    Write-Output "> go vet start"

    go vet ./...

    Write-Output "> go vet finished"

# reset path
Set-Location $path

# windows not allow run ps script:
# (admin start)Set-ExecutionPolicy RemoteSigned
