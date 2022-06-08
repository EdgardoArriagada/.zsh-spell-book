alias c='zsb_cat --paging=never'
alias le='zsb_cat --paging=always'
alias ls='lsd'
alias l='lsd -la'
alias removeSpanishChars='iconv -f utf8 -t ascii//TRANSLIT'
alias catenv='ls -a | rg env | xargs cat'
alias batenv='ls -a | rg env | xargs zsb_cat'

