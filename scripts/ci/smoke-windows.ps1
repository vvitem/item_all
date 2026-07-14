[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$Executable,

    [ValidateRange(1, 60)]
    [int]$ObservationSeconds = 10
)

$ErrorActionPreference = 'Stop'
$resolvedExecutable = (Resolve-Path -LiteralPath $Executable).Path
$process = Start-Process -FilePath $resolvedExecutable -PassThru

try {
    Start-Sleep -Seconds $ObservationSeconds
    $process.Refresh()
    if ($process.HasExited) {
        throw "ItemAll exited during the $ObservationSeconds second observation window with code $($process.ExitCode)."
    }
    Write-Host "ItemAll remained running for $ObservationSeconds seconds."
}
finally {
    $process.Refresh()
    if (-not $process.HasExited) {
        $cleanup = Start-Process -FilePath taskkill.exe -ArgumentList @('/PID', $process.Id, '/T', '/F') -Wait -PassThru -NoNewWindow
        if ($cleanup.ExitCode -ne 0) {
            throw "Failed to clean the ItemAll process tree; taskkill exit code $($cleanup.ExitCode)."
        }
    }
}
