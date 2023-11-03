# firstArg can be "--staged, --unmerged, etc"
${zsb}.recentGitFile() {
  local callback=$1
  local firstArg=$2

  if [[ "$firstArg" =~ "^--" ]]
    then eval "$callback ${@:3}"
    else eval "$callback $firstArg"
  fi
}

_${zsb}.recentGitFile() {
  local FIRST_ITEM_INDEX=3
  local usedCompletion=(${words[@]:2:$CURRENT-3})
  local firstItemUsed=${words[$FIRST_ITEM_INDEX]} # first item can be "--staged, --unmerged, etc or a file"
  local -a completionList

  # if we are completing the first item
  if [[ "$CURRENT" = "$FIRST_ITEM_INDEX" && "$firstItemUsed" =~ '^-' ]]; then
    for key in "${(@k)ZSB_GIT_FILETYPE_TO_REGEX}"; do
      if [[ "--${key}" =~ "^${firstItemUsed}" ]]; then
         completionList+=( --${key} )
      fi
    done

    _describe 'command' completionList
    return
  fi

  if [[ "$firstItemUsed" =~ '^--' ]]
    then completionList=( `${zsb}.getGitFiles ${firstItemUsed:2}` )
    else completionList=( `${zsb}.getGitFiles` )
  fi

  local newCompletion=( ${completionList:|usedCompletion} )

  _describe 'command' newCompletion
}

compdef _${zsb}.recentGitFile ${zsb}.recentGitFile

alias vr="${zsb}.recentGitFile v"
alias cr="${zsb}.recentGitFile c"
alias ccpr="${zsb}.recentGitFile ccp"
alias tigr="${zsb}.recentGitFile tig"
alias rmr="${zsb}.recentGitFile 'rm -rf'"
alias mvr="${zsb}.recentGitFile mv"
alias lr="${zsb}.recentGitFile l"
alias cpnr="${zsb}.recentGitFile cpn"
alias cdpr="${zsb}.recentGitFile cdp"
alias cdr="cdpr"
