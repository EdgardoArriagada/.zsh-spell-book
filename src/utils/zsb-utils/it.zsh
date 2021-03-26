# example: $ echo "Text that is $(it "italic") shows in italic"
it() echo "\e[3m$1\e[0m"
