# Sauvegarder les variables d'environnement actuelles
$originalGOOS = $env:GOOS
$originalGOARCH = $env:GOARCH
$originalCGO_ENABLED = $env:CGO_ENABLED

# Définir les variables d'environnement pour Windows
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "1"

try {
    go build -ldflags "-H=windowsgui" -o PoENotifier-windows-amd64.exe .
} finally {
    # Réinitialiser les variables d'environnement
    $env:GOOS = $originalGOOS
    $env:GOARCH = $originalGOARCH
    $env:CGO_ENABLED = $originalCGO_ENABLED
}
