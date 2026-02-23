#!/usr/bin/env pwsh

# Sem param() - leitura direta de variáveis de ambiente
$Version     = if ($env:TML_VERSION) { $env:TML_VERSION } else { "latest" }
$ResolvedHost = $env:TML_HOST
$ResolvedPort = $env:TML_PORT

$ArchitectureLabel = ""
if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
  $ArchitectureLabel = "x64"
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
  $ArchitectureLabel = "arm64"
} else {
  Write-Output "Falha na instalação: Windows $env:PROCESSOR_ARCHITECTURE não suportado.`n"
  return 1
}

$ErrorActionPreference = "Stop"

function Publish-EnvironmentVariables {
  if (-not ("Win32.NativeMethods" -as [Type])) {
    Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
  IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
  uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
  }
  $HwndBroadcast = [IntPtr] 0xffff
  $WmSettingChange = 0x1a
  $ResultPointer = [UIntPtr]::Zero
  [Win32.NativeMethods]::SendMessageTimeout(
    $HwndBroadcast, $WmSettingChange, [UIntPtr]::Zero,
    "Environment", 2, 5000, [ref] $ResultPointer
  ) | Out-Null
}

function Write-EnvironmentVariable {
  param([String]$KeyName, [String]$KeyValue)
  $RegistryKey = Get-Item -Path 'HKCU:'
  $EnvironmentKey = $RegistryKey.OpenSubKey('Environment', $true)
  if ($null -eq $KeyValue) {
    $EnvironmentKey.DeleteValue($KeyName)
  } else {
    $RegistryValueKind = if ($KeyValue.Contains('%')) { [Microsoft.Win32.RegistryValueKind]::ExpandString } else { [Microsoft.Win32.RegistryValueKind]::String }
    $EnvironmentKey.SetValue($KeyName, $KeyValue, $RegistryValueKind)
  }
  Publish-EnvironmentVariables
}

function Get-EnvironmentVariable {
  param([String]$KeyName)
  $RegistryKey = Get-Item -Path 'HKCU:'
  $EnvironmentKey = $RegistryKey.OpenSubKey('Environment')
  $EnvironmentKey.GetValue($KeyName, $null, [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames)
}

function Install-TMLSync {
  param(
    [string]$TargetVersion,
    [string]$HostIp,
    [string]$Port
  )

  $InstallationRoot = "${Home}\.tml-sync"
  $BinaryDirectory  = "${InstallationRoot}\bin"
  New-Item -ItemType Directory -Force -Path $BinaryDirectory | Out-Null

  $TargetFileName = "client-windows-$ArchitectureLabel.exe"
  $DownloadUrl    = "https://github.com/Ashu11-A/tModLoader-sync/releases/$(if ($TargetVersion -eq 'latest') { 'latest/download' } else { "download/$TargetVersion" })/$TargetFileName"
  $ExecutablePath = "${BinaryDirectory}\tml-sync.exe"

  Write-Output "Baixando tModLoader-sync ($TargetVersion) para windows/$ArchitectureLabel..."

  if (Get-Command "curl.exe" -ErrorAction SilentlyContinue) {
    curl.exe "-#SfLo" "$ExecutablePath" "$DownloadUrl"
  } else {
    Invoke-RestMethod -Uri $DownloadUrl -OutFile $ExecutablePath
  }

  if (-not (Test-Path $ExecutablePath)) {
    Write-Output "Falha na instalação - não foi possível baixar de $DownloadUrl"
    return 1
  }

  $CurrentPath = (Get-EnvironmentVariable -KeyName "Path") -split ';'
  if ($CurrentPath -notcontains $BinaryDirectory) {
    $CurrentPath += $BinaryDirectory
    Write-EnvironmentVariable -KeyName 'Path' -KeyValue ($CurrentPath -join ';')
    $env:PATH = $CurrentPath -join ';'
  }

  $ColorReset = [char]27 + "[0m"
  $ColorGreen = [char]27 + "[1;32m"
  Write-Output "${ColorGreen}tModLoader-sync instalado com sucesso!${ColorReset}"

  if ([string]::IsNullOrWhiteSpace($HostIp)) {
    $HostIp = Read-Host "Digite o IP do servidor"
  }

  if ([string]::IsNullOrWhiteSpace($Port)) {
    $Port = Read-Host "Digite a porta do servidor (ex: 25005)"
  }

  $FormattedPort = $Port.TrimStart(':')
  Write-Output "${ColorGreen}Iniciando tml-sync conectado a ${HostIp}:${FormattedPort}...${ColorReset}"

  & "$ExecutablePath" --host "$HostIp" --port "$FormattedPort"
}

Install-TMLSync -TargetVersion $Version -HostIp $ResolvedHost -Port $ResolvedPort