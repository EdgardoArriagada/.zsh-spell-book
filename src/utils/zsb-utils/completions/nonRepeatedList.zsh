# TODO Refactor to remove dup code
_${zsb}.nonRepeatedList() {
  local usedCompletion=( "${words[@]:1:$CURRENT-2}" )
  local completionList=( "$@" )

  local newCompletion=( ${completionList:|usedCompletion} )

  _describe 'command' newCompletion
}

# nonRepeatedList with Descriptions
_${zsb}.nonRepeatedListD() {
  local usedCompletion=( "${words[@]:1:$CURRENT-2}" )
  local completionList=( "$@" )

  declare -A tokenDictionary=( )
  local completionTokens=(  )
  for item in "${completionList[@]}"; do
    local token="${item%%:*}"
    tokenDictionary[$token]="$item"
    completionTokens+=( "$token" )
  done

  local newTokens=( ${completionTokens:|usedCompletion} )

  local newCompletion=(  )
  for token in "${newTokens[@]}"; do
    newCompletion+=( "${tokenDictionary[$token]}" )
  done

  _describe 'command' newCompletion
}

# nonRepeatedList with Command
_${zsb}.nonRepeatedListC() {
  local usedcompletion=( "${words[@]:1:$CURRENT-2}" )
  local completionlist=( $(eval "$1") )

  local newcompletion=( ${completionlist:|usedcompletion} )
  _describe 'command' newcompletion
}

