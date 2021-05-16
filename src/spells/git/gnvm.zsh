gnvm(){
  ${zsb}.warning "All your uncommited work will be lost"
  echo " "
  ${zsb}.prompt "Are you sure? [Y/n]"
  ${zsb}.confirmMenu && \
  git reset -q && git checkout . && git clean -fd
}

_${zsb}.nocompletion gnvm

