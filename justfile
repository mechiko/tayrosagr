shebang := 'pwsh.exe'
# Variables
exe_name := "tayrosagr"
mod_name := "tayrosagr"
ld_flags :="-H=windowsgui -s -w"
dist := ".dist"

default:
  just --list

win64:
    #!{{shebang}}
    $env:Path = "C:\Go\go.125\bin;C:\go\gcc\mingw64\bin;" + $env:Path
    $env:GOARCH = "amd64"
    $env:GOOS = "windows"
    $env:CGO_ENABLED = 1
    if (-Not (Test-Path go.mod)) {
      go mod init {{mod_name}}
    }
    go mod tidy -go 1.24 -v
    if(-Not $?) { exit }
    if (-Not (Test-Path "{{dist}}")) { New-Item -ItemType Directory -Force -Path "{{dist}}" | Out-Null }
    Remove-Item -Force -ErrorAction SilentlyContinue -LiteralPath "{{dist}}\{{exe_name}}.exe","{{dist}}\{{exe_name}}_64.exe"
    go build -ldflags="{{ld_flags}}" -o "{{dist}}\{{exe_name}}_64.exe" .
    if(-Not $?) { exit }
    upx --force-overwrite -o {{dist}}\{{exe_name}}.exe {{dist}}\{{exe_name}}_64.exe

win32:
    #!{{shebang}}
    $env:Path = "C:\Go\32\120\bin;C:\go\gcc\mingw32\bin;" + $env:Path
    $env:GOARCH = "386"
    $env:GOOS = "windows"
    $env:CGO_ENABLED = 1
    if (-Not (Test-Path go.mod)) {
      go mod init {{mod_name}}
    }
    go mod tidy -go 1.20 -v
    Set-Content -Path "go.mod" -Value (get-content -Path "go.mod" | Select-String -Pattern 'toolchain' -NotMatch)
    if(-Not $?) { exit }
    if (-Not (Test-Path "{{dist}}")) { New-Item -ItemType Directory -Force -Path "{{dist}}" | Out-Null }
    Remove-Item -Force -ErrorAction SilentlyContinue -LiteralPath "{{dist}}\{{exe_name}}.exe","{{dist}}\{{exe_name}}_32.exe"
    go build -ldflags="{{ld_flags}}" -o "{{dist}}\{{exe_name}}_32.exe" .
    if(-Not $?) { exit }
    upx --force-overwrite -o {{dist}}\{{exe_name}}32.exe {{dist}}\{{exe_name}}_32.exe
