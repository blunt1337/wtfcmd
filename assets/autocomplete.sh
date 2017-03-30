# Should be used with:
# eval "$(CMDNAME --autocomplete setup)"

# Call our program to autocomplete
_CMDNAME_completion_loader() {
	local string="$(CMDPATH --autocomplete $COMP_POINT "$COMP_LINE" "${COMP_WORDS[@]}")"
	local index=0
	local sep="SEPARATOR"
	
	# Empty
	if [ -z "$string" ]; then
		COMPREPLY=()
        return 0
	fi
	
	# Split into COMPREPLY
	while [ "$string" != "${string#*$sep}" ]; do
		COMPREPLY[$index]="${string%%$sep*}"
		index=$index+1
		string="${string#*$sep}"
	done
	COMPREPLY[$index]="$string"
}
complete -F _CMDNAME_completion_loader -o default CMDNAME