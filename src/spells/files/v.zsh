v() {
  local -r file="$1"

  case "$file" in
    *.xlsx|*.csv|*.odt) libreoffice --calc "$file" ;;

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
    |*.svg) sxiv -f "$file" ;;

    *) eval "nnvim ${file}" && ${zsb}.isGitRepo && ${zsb}.gitStatus ;;
  esac
}

