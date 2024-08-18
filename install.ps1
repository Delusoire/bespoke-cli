[CmdletBinding()]
param (
	[Parameter(
		Mandatory = $false,
		Position = 0,
		HelpMessage = "Specify the version of Spicetify to install (format: 1.2.3). If not specified, the latest version will be installed."
	)]
	[string]$v,
	[Parameter(
		Mandatory = $false,
		HelpMessage = "Specify the path to the Spicetify folder."
	)]
	[string]$spicetifyFolderPath = "$env:LOCALAPPDATA\Spicetify",
	[Parameter(
		Mandatory = $false,
		HelpMessage = "Install Spicetify in portable mode. Storing the configuration within <spicetifyFolderPath>\config."
	)]
	[switch]$portable = $false,
	[Parameter(
		Mandatory = $false,
		HelpMessage = "Skip installing the binary."
	)]
	[switch]$skipBinary = $false,
	[Parameter(
		Mandatory = $false,
		HelpMessage = "Skip initializing the daemon."
	)]
	[switch]$skipDaemon = $false,
	[Parameter(
		Mandatory = $false,
		HelpMessage = "Skip registering the URI scheme."
	)]
	[switch]$skipURIScheme = $false
)

$ErrorActionPreference = 'Stop'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

#region Variables
$cliOwnerRepo = "Delusoire/bespoke-cli"

$spicetifyPortableBinaryPath = "$spicetifyFolderPath\bin"
$spicetifyPortableExecutablePath = "$spicetifyPortableBinaryPath\spicetify.exe"
#endregion Variables

#region Functions
function Write-Ok {
	[CmdletBinding()]
	param ()
	process {
		Write-Host -Object ' > OK' -ForegroundColor 'Green'
	}
}

function Write-Error {
	[CmdletBinding()]
	param ()
	process {
		Write-Host -Object ' > ERROR' -ForegroundColor 'Red'
	}
}

function Test-Admin {
	[CmdletBinding()]
	param ()
	begin {
		Write-Host -Object "Checking if the script is not being run as administrator..." -NoNewline
	}
	process {
		$currentUser = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
		-not $currentUser.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
	}
}

function Test-PowerShellVersion {
	[CmdletBinding()]
	param ()
	begin {
		$PSMinVersion = [version]'5.1'
	}
	process {
		Write-Host -Object 'Checking if your PowerShell version is compatible...' -NoNewline
		$PSVersionTable.PSVersion -ge $PSMinVersion
	}
}

function Add-Folder {
	[CmdletBinding()]
	param ()
	begin {
		Write-Host -Object 'Creating Spicetify folder...' -NoNewline
	}
	process {
		if (Test-Path -Path $spicetifyFolderPath) {
			Move-Item -Path $spicetifyFolderPath -Destination "$spicetifyFolderPath.old" -Force
		}
		New-Item -Path $spicetifyFolderPath -ItemType 'Directory' -Force
	}
	end {
		Write-Ok
	}
}

function Receive-Binary {
	[CmdletBinding()]
	param ()
	begin {
		if ($env:PROCESSOR_ARCHITECTURE -eq 'AMD64') {
			$architecture = 'amd64'
		}
		elseif ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') {
			$architecture = 'arm64'
		}
		else {
			$architecture = '386'
		}
		if ($v) {
			if ($v -match '^\d+\.\d+\.\d+$') {
				$targetVersion = "v$v"
			}
			else {
				Write-Warning -Message "You have spicefied an invalid Spicetify version: $v `nThe version must be in the following format: 1.2.3"
				Pause
				exit
			}
		}
		else {
			Write-Host -Object 'Fetching the latest Spicetify version...' -NoNewline
			$latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/$cliOwnerRepo/releases/latest"
			$targetVersion = $latestRelease.tag_name
			Write-Ok
		}
		$binaryPath = [System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), "spicetify.exe")
	}
	process {
		Write-Host -Object "Downloading Spicetify $targetVersion..." -NoNewline
		$Parameters = @{
			Uri            = "https://github.com/$cliOwnerRepo/releases/download/$targetVersion/bespoke-cli-$v-windows-$architecture.exe"
			UseBasicParsin = $true
			OutFile        = $binaryPath
		}
		Invoke-WebRequest @Parameters
		Write-Ok
	}
	end {
		$binaryPath
	}
}

function Add-BinToPath {
	[CmdletBinding()]
	param ()
	begin {
		Write-Host -Object 'Making Spicetify available in the PATH...' -NoNewline
		$user = [EnvironmentVariableTarget]::User
		$path = [Environment]::GetEnvironmentVariable('PATH', $user)
	}
	process {
		if ($path -notlike "*$spicetifyPortableBinaryPath*") {
			$path = "$path;$spicetifyPortableBinaryPath"
		}
	}
	end {
		[Environment]::SetEnvironmentVariable('PATH', $path, $user)
		$env:PATH = $path
		Write-Ok
	}
}

function Install-Binary {
	[CmdletBinding()]
	param ()
	begin {
		Write-Host -Object 'Installing Spicetify...'
		Add-Folder
	}
	process {
		$downloadedBinaryPath = Receive-Binary
		Write-Host -Object 'Extracting Spicetify...' -NoNewline
		New-Item -Path $spicetifyPortableBinaryPath -ItemType 'Directory' -Force
		Move-Item -Path $downloadedBinaryPath -DestinationPath $spicetifyPortableBinaryPath -Force
		Write-Ok
		Add-BinToPath
		if ($portable) {
			Write-Host -Object 'Creating Spicetify portable config folder...' -NoNewline
			$spicetifyPortableConfigPath = "$spicetifyFolderPath\config"
			New-Item -Path $spicetifyPortableConfigPath -ItemType 'Directory' -Force
			Write-Ok
		}
		& $spicetifyPortableExecutablePath 'init'
	}
	end {
		Write-Host -Object 'Spicetify was successfully installed!' -ForegroundColor 'Green'
	}
}

function Initialize-Daemon {
	[CmdletBinding()]
	param ()
	begin {
		Write-Host -Object 'Creating Spicetify daemon task...' -NoNewline
	}
	process {
		$initTask = {
			$taskName = "Spicetify daemon"
			$description = "Launches Spicetify daemon at startup"
			$command = $spicetifyPortableExecutablePath
			$arguments = "daemon"

			$action = New-ScheduledTaskAction -Execute $command -Argument $arguments
			$trigger = New-ScheduledTaskTrigger -AtStartup
			$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -ExecutionTimeLimit (New-TimeSpan -Seconds 0)
			$principal = New-ScheduledTaskPrincipal -UserId $env:USERNAME -LogonType S4U -RunLevel Limited

			$task = New-ScheduledTask -Action $action -Trigger $trigger -Settings $settings -Principal $principal -Description $description

			Register-ScheduledTask -TaskName $taskName -InputObject $task
			Start-ScheduledTask -TaskName $taskName
		}

		if (Test-Admin) {
			Invoke-Expression $initTask
		}
		else {
			Write-Host -Object 'Running the task as administrator...' -NoNewline
			$tempFile = [System.IO.Path]::GetTempFileName()
			$tempFile += ".ps1"
			$initTask | Out-File -FilePath $tempFile -Encoding UTF8
			Start-Process powershell "-NoProfile -ExecutionPolicy Bypass -File `"$tempFile`"" -PassThru -Verb RunAs -WindowStyle Hidden -Wait
			Write-Ok
		}
	}
	end {
		Write-Host -Object 'Deamon task was successfully created!' -ForegroundColor 'Green'
	}
}

function Register-URIScheme {
	[CmdletBinding()]
	param ()
	begin {
		Write-Host -Object 'Registering Spicetify URI scheme...' -NoNewline
	}
	process {
		$scheme = "spicetify"
		$command = "`"$spicetifyPortableExecutablePath`" protocol `"%1`""

		$K = New-Item -Path "HKCU:\Software\Classes\$scheme" -Force
		$K.SetValue("", "URL:$scheme Protocol", [Microsoft.Win32.RegistryValueKind]::String)
		$K = $K.CreateSubKey("shell\open\command")
		$K.SetValue("", "$command", [Microsoft.Win32.RegistryValueKind]::String)
	}
	end {
		Write-Host -Object 'URI scheme was successfully registered!' -ForegroundColor 'Green'
	}
}
#endregion Functions

#region Main
#region Checks
if (-not (Test-PowerShellVersion)) {
	Write-Error
	Write-Warning -Message 'PowerShell 5.1 or higher is required to run this script'
	Write-Warning -Message "You are running PowerShell $($PSVersionTable.PSVersion)"
	Write-Host -Object 'PowerShell 5.1 install guide:'
	Write-Host -Object 'https://learn.microsoft.com/skypeforbusiness/set-up-your-computer-for-windows-powershell/download-and-install-windows-powershell-5-1'
	Write-Host -Object 'PowerShell 7 install guide:'
	Write-Host -Object 'https://learn.microsoft.com/powershell/scripting/install/installing-powershell-on-windows'
	Pause
	exit
}
else {
	Write-Ok
}
if (-not (Test-Admin)) {
	Write-Error
	Write-Warning -Message "The script was run as administrator. This can result in problems with the installation process or unexpected behavior. Do not continue if you do not know what you are doing."
	$Host.UI.RawUI.Flushinputbuffer()
	$choices = [System.Management.Automation.Host.ChoiceDescription[]] @(
    (New-Object System.Management.Automation.Host.ChoiceDescription '&Yes', 'Abort installation.'),
    (New-Object System.Management.Automation.Host.ChoiceDescription '&No', 'Resume installation.')
	)
	$choice = $Host.UI.PromptForChoice('', 'Do you want to abort the installation process?', $choices, 0)
	if ($choice -eq 0) {
		Write-Host -Object 'Spicetify installation aborted' -ForegroundColor 'Yellow'
		Pause
		exit
	}
}
else {
	Write-Ok
}
#endregion Checks

#region Spicetify
if (-not $skipBinary) {
	Install-Binary
}
if (-not $skipDaemon) {
	Initialize-Daemon
}
if (-not $skipURIScheme) {
	Register-URIScheme
}

Write-Host -Object "`nRun" -NoNewline
Write-Host -Object ' spicetify -h ' -NoNewline -ForegroundColor 'Cyan'
Write-Host -Object 'to get started'
#endregion Spicetify
#endregion Main
