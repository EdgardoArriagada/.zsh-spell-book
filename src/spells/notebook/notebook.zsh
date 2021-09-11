notebook() {
  ${zsb}.validate "ZSB_NOTEBOOK_DIR"
  ${zsb}.validate "ZSB_NOTEBOOK_CHAPTER"
  ${zsb}.validate "ZSB_NOTEBOOK_PAGE"

  local -r currentPage="${ZSB_NOTEBOOK_DIR}/${ZSB_NOTEBOOK_CHAPTER}/${ZSB_NOTEBOOK_PAGE}"

  # Create parent folder
  mkdir -p $(dirname $currentPage)

  eval "nnvim ${currentPage}"
}
