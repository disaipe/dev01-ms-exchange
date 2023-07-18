$credentials = Get-Content ./credentials
$parts = $credentials.Split("|")

$_host = $parts[0]
$_login = $parts[1]
$_password = $parts[2] | ConvertTo-SecureString

$cred = New-Object System.Management.Automation.PSCredential($_login, $_password)
$session = New-PSSession -ComputerName $_host -Credential $cred -Authentication CredSSP

$OutputEncoding = [Console]::OutputEncoding = New-Object System.Text.Utf8Encoding

if ($session) {
    Enter-PSSession -Session $session

    Invoke-Command -Session $session -ScriptBlock {
        Add-PSSnapin Microsoft.Exchange.Management.PowerShell.SnapIn
        $databases = Get-MailboxDatabase -Status | Where-Object { $_.Mounted -eq 'True' }
        $results = @()
        foreach ($database in $databases) {
            $items = Get-MailboxStatistics -Database $database
            foreach ($item in $items) {
                $prop = @{
                    Id=$item.Identity.MailboxGuid.Guid
                    DisplayName=$item.DisplayName
                    TotalItemSize=$item.TotalItemSize.Value.ToBytes()
                    TotalItemCount=$item.ItemCount
                }
                $results += $prop
           }
        }
        return $results | ConvertTo-Json -Compress
    }

    Remove-PSSession -Session $session
}