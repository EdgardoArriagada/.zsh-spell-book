# Vim Directory, a "vim that change to the directory of the file opend"
# the following effectis will occur depending on arguments:
# 1.- if you give a full file path, it should cd to that file's directory and list it (after file is closed)
# 2.- if you give a directory but no file, it should just cd to the directory given and list it
# 3.- if no file or directory is given, it should just cd and ls
# 4.- if you open a file that is inside a git repo, it should do item 1.- and it should call git status as well (after file is closed)
# 5.- if $1 is not a valid file or directory, it should create a new file out of it, calling git status after closing it

vimd() {
  local this="$0_$(${zsb}.timeId)"
  {
    ${this}.main() {
      local USER_INPUT=$1

      if [ -z "$USER_INPUT" ]; then
        builtin cd && ls
        return 0
      fi

      if [ -d "$USER_INPUT" ]; then
        builtin cd "$USER_INPUT" && ls
        return 0
      fi

      ${this}.handleFullFilePath "$USER_INPUT"
      return 0
    }

    ${this}.handleFullFilePath() {
      local fullFilePath=$1
      local fileDirectoryName=$(dirname ${fullFilePath})
      vim "$fullFilePath" && cds "$fileDirectoryName"
      return 0
    }

    ${this}.main "$@"

  } always {
    unfunction -m "${this}.*"
  }
}

alias vid="vimd"
alias vd="vimd"
