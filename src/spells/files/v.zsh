if (( $ZSB_MACOS )); then
  ${zsb}.v.openCsv() open $1
  ${zsb}.v.openImg() open $1
  ${zsb}.v.openVideo() open $1
else
  ${zsb}.v.openCsv() libreoffice --calc $1
  ${zsb}.v.openImg() sxiv -f $1
  ${zsb}.v.openVideo() vlc $1
fi

${zsb}.v() {
  case ${1:e:l} in
    xlsx|csv|odt) ${0}.openCsv $1 ;;

    mp4|mkv|avi|mov|webm|flv|wmv) ${0}.openVideo $1 ;;

    pdf|jpg|png|gif|webp|tiff|psd|raw|bmp|heif|jpeg|svg) ${0}.openImg $1 ;;

    *) nvim $1 && ${zsb}.tryGitStatus ;;
  esac
}

alias v='noglob ${zsb}.v'

hisIgnore v

