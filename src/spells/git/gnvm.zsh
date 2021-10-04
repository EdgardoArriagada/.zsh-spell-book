gnvm(){
  ${zsb}.fullPrompt "All your uncommited work will be lost"

  ${zsb}.confirmMenu && \
  git reset -q && git checkout . && git clean -fd
}

_${zsb}.nocompletion gnvm

