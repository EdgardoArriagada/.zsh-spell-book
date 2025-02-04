${zsb}.cpdate_generic() { zsb_clipcopy $1 && ${zsb}.success "`hl $1` as `hl $2` copied to clipboard."; }

alias cpdate="${zsb}.cpdate_generic `date +%Y-%m-%d` YYYY-MM-DD"
alias cpdate2="${zsb}.cpdate_generic `date +%d-%m-%Y` DD-MM-YYYY"

