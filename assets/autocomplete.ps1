# Autocomplete function
$completion_CMDNAME = {
	param($commandName, $commandAst, $cursorPosition)

	# Params for CMDNAME --autocomplete
	$cmdline = $commandAst.ToString()
	$words = @()
	foreach ($word in $commandAst.CommandElements) {
		$words += $word.ToString()
	}

	# Results
	$output = CMDNAME --autocomplete $cursorPosition "$cmdline" $words | Out-String
	$output = $output -replace "\r"
	$output = $output -split "SEPARATOR"
	foreach ($word in $output) {
		if ($word -ne '') {
			New-Object System.Management.Automation.CompletionResult $word, $word, 'ParameterValue', $word
		}
	}
}

# Register the TabExpension2 function
if (-not $global:options) { $global:options = @{CustomArgumentCompleters = @{};NativeArgumentCompleters = @{}}}
$global:options['NativeArgumentCompleters']['CMDNAME'] = $completion_CMDNAME
$function:tabexpansion2 = $function:tabexpansion2 -replace 'End\r\n{','End { if ($null -ne $options) { $options += $global:options} else {$options = $global:options}'