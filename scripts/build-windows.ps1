$originalGOOS = $env:GOOS
$originalGOARCH = $env:GOARCH
# $originalCGO_ENABLED = $env:CGO_ENABLED

# $env:GOOS = "windows"
# $env:GOARCH = "amd64"
# $env:CGO_ENABLED = "1"

try {
    go build -ldflags -H=windowsgui -o PoENotifier-windows-amd64.exe .
} finally {
    $env:GOOS = $originalGOOS
    $env:GOARCH = $originalGOARCH
    $env:CGO_ENABLED = $originalCGO_ENABLED
}
