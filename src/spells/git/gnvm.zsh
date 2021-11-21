gnvm(){
  ${zsb}.confirmMenu.withItems "All your uncommited work will be lost"

  git reset -q && git checkout . && git clean -fd
}

_${zsb}.nocompletion gnvm

