visualizeAll() {
  local -r extension=${1:?'You must provide a file extension'}

  for file in ./*.${extension}; do
    ${zsb}.info "Visualizing ${file}"
    eval "v '${file}'"
  done
}

# TODO: complete with all file extensions in current directory

