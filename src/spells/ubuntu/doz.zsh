doz() {
  [[ -z "$1" ]] && ${zsb}.throw 'Usage: doz <file_to_compress>'

  local outputFile="${1}.zip"
  [[ -f "$outputFile" ]] && ${zsb}.throw "$(hl "$outputFile") already exists."

  if [[ -f "$1" ]]; then
    zip "$outputFile" "$1" && ${zsb}.success "Created $(hl "$outputFile")"
    return $?
  fi

  if [[ -d "$1" ]]; then
    zip -r "$outputFile" "$1" && ${zsb}.success "Created $(hl "$outputFile")"
    return $?
  fi

  ${zsb}.throw "$(hl "$1") is not a valid file or directory."
}
