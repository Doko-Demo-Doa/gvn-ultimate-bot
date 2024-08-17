New-Item -Name "out" -ItemType "directory" -Force
$curr = Get-Location
$tempdir = "out"
$target = "$curr\$tempdir"

$env:GOTMPDIR = $target
$env:GO111MODULE = "on"
go build main.go