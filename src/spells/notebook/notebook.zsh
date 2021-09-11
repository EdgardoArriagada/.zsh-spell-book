notebook() {
  ${zsb}.validate "ZSB_NOTEBOOK_DIR"
  ${zsb}.validate "ZSB_NOTEBOOK_CHAPTER"
  ${zsb}.validate "ZSB_NOTEBOOK_PAGE"

  eval "nnvim ${ZSB_NOTEBOOK_DIR}/${ZSB_NOTEBOOK_CHAPTER}/${ZSB_NOTEBOOK_PAGE}"
}
