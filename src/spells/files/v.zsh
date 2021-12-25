v() {
  local -r file="$1"

  case "$file" in
    *.xlsx|*.csv|*.odt) (( $ZSB_MACOS )) || libreoffice --calc "$file" ;;

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
    |*.svg) (( $ZSB_MACOS )) && open "$file" || sxiv -f "$file" ;;

    *) eval "nnvim ${file}" && ${zsb}.isGitRepo && ${zsb}.gitStatus ;;
  esac
}

