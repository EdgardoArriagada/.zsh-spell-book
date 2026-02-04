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
  local file="$1"
  local line_num=""

  # Check if input contains :line_number pattern
  if [[ "$file" =~ ^(.+):([0-9]+)$ ]]; then
    file="${match[1]}"
    line_num="${match[2]}"
  fi

  case ${file:e:l} in
    xlsx|odt) ${0}.openCsv "$file" ;;

    mp4|mkv|avi|mov|webm|flv|wmv) ${0}.openVideo "$file" ;;

    pdf|jpg|png|gif|webp|tiff|psd|raw|bmp|heif|jpeg|svg) ${0}.openImg "$file" ;;

    *)
      if [[ -n "$line_num" ]]
        then nvim +"$line_num" "$file" && ${zsb}.tryGitStatus
        else nvim "$file" && ${zsb}.tryGitStatus
      fi
      ;;
  esac
}

alias v='noglob ${zsb}.v'

hisIgnore v

