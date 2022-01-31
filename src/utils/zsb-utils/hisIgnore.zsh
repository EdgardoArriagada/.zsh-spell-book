declare -a ZSB_HISTORY_IGNORE

hisIgnore() ZSB_HISTORY_IGNORE+=( "$1" )

# the following line should be applied as an automatic call
# export HISTORY_IGNORE="(${(j:|:)ZSB_HISTORY_IGNORE})"
