${zsb}.ve() {
  if [[ "$1" =~ '/' && ! -a $1 ]]; then
    if ! ${zsb}.isGitRepo || [[ "$1" =~ '^(/|~|\.\.)' ]]
      then ${zsb}.throw 'Unsafe to perform roman expansion.'
    fi

    if [[ "${1: -1}" == "/" ]]
      then mkdir -p $1
      else mkdir -p `dirname $1`
    fi

    nvim $1 && ${zsb}.tryGitStatus

    return 0
  fi

  nvim $1 && ${zsb}.tryGitStatus
}

alias ve='noglob ${zsb}.ve'

hisIgnore ve

