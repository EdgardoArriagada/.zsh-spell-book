# Disable suspend the output "feature"
stty stop undef # disables C-s
stty start '^-' stop '^-' # disables C-q

# Editor for `fc` command
FCEDIT=nvim
