# print without escaping special characters nor parcing args
# pass args in quotes
plainPrint() echo -E - "$1"

alias puts='plainPrint'

