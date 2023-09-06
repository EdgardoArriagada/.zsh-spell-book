if (( $ZSB_MACOS )); then
  ${zsb}.v.openCsv() open "${1}"
  ${zsb}.v.openImg() open "${1}"
  ${zsb}.v.openVideo() open "${1}"
else
  ${zsb}.v.openCsv() libreoffice --calc "${1}"
  ${zsb}.v.openImg() sxiv -f "${1}"
  ${zsb}.v.openVideo() vlc "${1}"
fi

v() {
  local file=$1

  case ${file:l} in
    *.xlsx|*.csv|*.odt) ${zsb}.v.openCsv ${file} ;;

    *.mp4|*.mkv|*.avi|*.mov|*.webm|*.flv|*.wmv) ${zsb}.v.openVideo ${file} ;;

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
    |*.svg) ${zsb}.v.openImg ${file} ;;

    *) nvim ${file} && ${zsb}.isGitRepo && ${zsb}.gitStatus ;;
  esac
}

hisIgnore v

