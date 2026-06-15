$link = "https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/EquilotlCli.exe"

$outfile = "$env:TEMP\EquilotlCli.exe"

Write-Output "Downloading installer to $outfile"

Invoke-WebRequest -Uri "$link" -OutFile "$outfile"

Write-Output ""

Start-Process -Wait -NoNewWindow -FilePath "$outfile"

# Cleanup
Remove-Item -Force "$outfile"
