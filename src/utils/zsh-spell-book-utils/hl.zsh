# example: $ echo "Text that is $(hl "highlighted") shows in color"
hl() echo "${ZSB_SHL}$1${ZSB_EHL}"
