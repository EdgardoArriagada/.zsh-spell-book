# replace words in the last history entry with
# arg(s) given to this function.

# for better experience, excecute this function beginning with a space
# so is not keept in history and you can keep replacing the 1st
# word in the last history entry again and again

again() {
  local substractions="$(($1 + 1))"
  shift 1

  local lastHistoryEntryArguments=$(history -n -1 | cut -d " " -f${substractions}-)

  local newCommand="${@} ${lastHistoryEntryArguments}"

  if [[ ! "$(echo $ZSH_VERSION)" ]]; then
    echo "${ZSB_INFO} Executing \"${newCommand}\""
    eval "$newCommand"
    return 0
  fi

  print -z "$newCommand"

}

alias ag0="again 0"
alias ag="again 1"
alias ag1="again 1"
alias ag2="again 2"
alias ag3="again 3"
alias ag4="again 4"
