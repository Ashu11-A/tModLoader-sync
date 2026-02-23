#!/usr/bin/env pwsh
param(
  [String]$Version = "latest",
  [String]$hostAddress = "",
  [String]$serverPort = ""
)

if ($hostAddress -eq "" -and $null -ne $global:hostAddress) { $hostAddress = $global:hostAddress }
if ($serverPort -eq "" -and $null -ne $global:serverPort) { $serverPort = $global:serverPort }

$ArchLabel = ""
if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
  $ArchLabel = "x64"
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
  $ArchLabel = "arm64"
} else {
  Write-Output "Installation failed: Windows $env:PROCESSOR_ARCHITECTURE not supported.`n"
  return 1
}

$ErrorActionPreference = "Stop"

function Publish-Env {
  if (-not ("Win32.NativeMethods" -as [Type])) {
    Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
  IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
  uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
  }
  $HWND_BROADCAST = [IntPtr] 0xffff
  $WM_SETTINGCHANGE = 0x1a
  $result = [UIntPtr]::Zero
  [Win32.NativeMethods]::SendMessageTimeout($HWND_BROADCAST,
    $WM_SETTINGCHANGE,
    [UIntPtr]::Zero,
    "Environment",
    2,
    5000,
    [ref] $result
  ) | Out-Null
}

function Write-Env {
  param([String]$Key, [String]$Value)
  $RegisterKey = Get-Item -Path 'HKCU:'
  $EnvRegisterKey = $RegisterKey.OpenSubKey('Environment', $true)
  if ($null -eq $Value) {
    $EnvRegisterKey.DeleteValue($Key)
  } else {
    $RegistryValueKind = if ($Value.Contains('%')) { [Microsoft.Win32.RegistryValueKind]::ExpandString } else { [Microsoft.Win32.RegistryValueKind]::String }
    $EnvRegisterKey.SetValue($Key, $Value, $RegistryValueKind)
  }
  Publish-Env
}

function Get-Env {
  param([String] $Key)
  $RegisterKey = Get-Item -Path 'HKCU:'
  $EnvRegisterKey = $RegisterKey.OpenSubKey('Environment')
  $EnvRegisterKey.GetValue($Key, $null, [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames)
}

function Install-TMLSync {
  param([string]$Version)

  $InstallRoot = "${Home}\.tml-sync"
  $BinDir = mkdir -Force "${InstallRoot}\bin"
  
  $Target = "client-windows-$ArchLabel.exe"
  $URL = "https://github.com/Ashu11-A/tModLoader-sync/releases/$(if ($Version -eq "latest") { "latest/download" } else { "download/$Version" })/$Target"

  $ExePath = "${BinDir}\tml-sync.exe"

  Write-Output "Downloading tModLoader-sync ($Version) for windows/$ArchLabel..."
  
  if (Get-Command "curl.exe" -ErrorAction SilentlyContinue) {
    curl.exe "-#SfLo" "$ExePath" "$URL"
  } else {
    Invoke-RestMethod -Uri $URL -OutFile $ExePath
  }

  if (!(Test-Path $ExePath)) {
    Write-Output "Installation failed - could not download $URL"
    return 1
  }

  $Path = (Get-Env -Key "Path") -split ';'
  if ($Path -notcontains $BinDir) {
    $Path += $BinDir
    Write-Env -Key 'Path' -Value ($Path -join ';')
    $env:PATH = $Path -join ';'
  }

  $C_RESET = [char]27 + "[0m"
  $C_GREEN = [char]27 + "[1;32m"
  Write-Output "${C_GREEN}tModLoader-sync installed successfully!${C_RESET}"
  
  if ($hostAddress -ne "" -and $serverPort -ne "") {
    $portVal = $serverPort.TrimStart(':')
    Write-Output "${C_GREEN}Starting tml-sync connected to ${hostAddress}:${portVal}...${C_RESET}"
    & "$ExePath" --host "$hostAddress" --port "$portVal"
  } else {
    Write-Output "To start syncing, run: tml-sync --host <IP> --port <PORT>"
  }
}

Install-TMLSync -Version $Version