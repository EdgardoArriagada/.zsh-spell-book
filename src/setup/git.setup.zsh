if git --version >/dev/null 2>&1; then
  # configuration goes to ~/.gitconfig

  git config --global core.editor "vim"
  git config --global --replace-all core.excludesfile ${ZSB_DIR}/src/dotFiles/git/global.gitignore

  git config --global pager.branch false # print branch info instead of "less" it
  git config --global pager.diff delta
  git config --global pager.log delta
  git config --global pager.reflog delta
  git config --global pager.show delta
  git config --global include.path ${ZSB_DIR}/src/dotFiles/git/themes.gitconfig
  git config --global interactive.diffFilter 'delta --color-only'
  git config --global delta.features calochortus-lyallii
fi

