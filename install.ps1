$ErrorActionPreference = 'Stop'
$architecture = $env:PROCESSOR_ARCHITECTURE
if ($env:PROCESSOR_ARCHITEW6432) {
    $architecture = $env:PROCESSOR_ARCHITEW6432
}
$binary = if ($architecture -eq 'ARM64') { 'IllegalotlCli-arm64.exe' } else { 'IllegalotlCli.exe' }
$link = "https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/$binary"

$outfile = Join-Path $env:TEMP ("IllegalcordInstaller-" + [guid]::NewGuid() + '.exe')

Write-Output "Downloading installer to $outfile"

try {
    Invoke-WebRequest -Uri $link -OutFile $outfile
    Write-Output ""
    $installer = Start-Process -Wait -NoNewWindow -PassThru -FilePath $outfile
    if ($installer.ExitCode -ne 0) {
        throw "Illegalcord Installer exited with code $($installer.ExitCode)"
    }
} finally {
    if (Test-Path -LiteralPath $outfile) {
        Remove-Item -Force -LiteralPath $outfile
    }
}
