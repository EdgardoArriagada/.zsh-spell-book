${zsb}.extensionOf() {
  printf "$1" | rev | cut -d"." -f1 | rev
}

