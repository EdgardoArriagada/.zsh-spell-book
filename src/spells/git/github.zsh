github.getNormalUrlFromSshUrl() {
  local input="${1:4}"
  local formattedInput=`printf "$input" | sed 's/:/\//g' | sed 's/\.git$//g'`
  printf "https://${formattedInput}"
}

github() {
  local remoteUrl=`git config --get remote.origin.url`

  ${zsb}.isUrl "$remoteUrl" && zsb_open "$remoteUrl" && return $?

  zsb_open `${0}.getNormalUrlFromSshUrl "$remoteUrl"`
}

hisIgnore github

_${zsb}.nocompletion github
