_${zsb}.nvm() {
  local compList
  case "$CURRENT" in
    2)
      compList=(
        'ls-remote:List remote node version'
        'install:<version> Installs given version'
        'use:<version> Use given version during terminal session'
        'alias:<tab>'
        'uninstall:<version> uninstall given version'
        'ls:List installed versions'
      ) ;;
    3)
      local -A subCmds
      subCmds[alias]="default:\ Set\ default\ node\ version other:\ Set\ custom\ alias"
      compList=( "${(z)subCmds[${words[CURRENT - 1]}]}" )
  esac


  _describe 'command' compList
}

compdef _${zsb}.nvm nvm

