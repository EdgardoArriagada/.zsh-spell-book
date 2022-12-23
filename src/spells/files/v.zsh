if (( $ZSB_MACOS )); then
  ${zsb}.v.openCsv() open "${1}"
  ${zsb}.v.openImg() open "${1}"
else
  ${zsb}.v.openCsv() libreoffice --calc "${1}"
  ${zsb}.v.openImg() sxiv -f "${1}"
fi

v() {
  local -r file="$1"

  case "$file" in
    *.xlsx|*.csv|*.odt) ${zsb}.v.openCsv "$file" ;;

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
    |*.svg) ${zsb}.v.openImg "$file" ;;

    *) nvim ${file} && ${zsb}.isGitRepo && ${zsb}.gitStatus ;;
  esac
}

hisIgnore v

