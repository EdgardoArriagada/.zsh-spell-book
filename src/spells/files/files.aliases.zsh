(( $ZSB_MACOS )) && alias batcat='bat'
alias c='batcat --paging=never'
alias le='batcat --paging=always'
alias ls='lsd'
alias removeSpanishChars='iconv -f utf8 -t ascii//TRANSLIT'
alias catenv='ls -a | rg env | xargs cat'
alias batcatenv='ls -a | rg env | xargs batcat'

