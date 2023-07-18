$targetHost = Read-Host "Enter target host with installed 'Microsoft.Exchange.Management.PowerShell.SnapIn'"

$credential = Get-Credential
$encryptedPass = $credential.Password | ConvertFrom-SecureString

"{0}|{1}|{2}" -f $targetHost, $credential.Username, $encryptedPass | Set-Content ./credentials