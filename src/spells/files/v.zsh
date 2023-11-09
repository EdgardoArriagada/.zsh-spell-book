if (( $ZSB_MACOS )); then
  ${zsb}.v.openCsv() open $1
  ${zsb}.v.openImg() open $1
  ${zsb}.v.openVideo() open $1
else
  ${zsb}.v.openCsv() libreoffice --calc $1
  ${zsb}.v.openImg() sxiv -f $1
  ${zsb}.v.openVideo() vlc $1
fi

${zsb}.visualize() {
  case ${1:l} in
    *.xlsx|*.csv|*.odt) ${zsb}.v.openCsv $1 ;;

    *.mp4|*.mkv|*.avi|*.mov|*.webm|*.flv|*.wmv) ${zsb}.v.openVideo $1 ;;

    *.pdf \
    |*.jpg \
    |*.png \
    |*.gif \
    |*.webp \
    |*.tiff \
    |*.psd \
    |*.raw \
    |*.bmp \
    |*.heif \
    |*.jpeg \
    |*.svg) ${zsb}.v.openImg $1 ;;

    *) nvim $1 && ${zsb}.isGitRepo && ${zsb}.gitStatus ;;
  esac
}

alias v='noglob ${zsb}.visualize'

hisIgnore v

