New-Item -Name "out" -ItemType "directory"
$curr = Get-Location
$tempdir = "out"
$target = "$curr\$tempdir"

$env:GOTMPDIR = $target
go run main.go