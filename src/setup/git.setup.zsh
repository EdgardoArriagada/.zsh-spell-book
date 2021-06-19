if git --version >/dev/null 2>&1; then
  # configuration goes to ~/.gitconfig

  git config --global core.editor "vim"
  git config --global pager.branch false # print branch info instead of "less" it
  git config --global --replace-all core.excludesfile ${ZSB_DIR}/src/dotFiles/git/global.gitignore

  # Highlight changed words in git diff (dpkg -L git | grep diff-highlight)
  git config --global pager.log "perl /usr/share/doc/git/contrib/diff-highlight/diff-highlight | less"
  git config --global pager.show "perl /usr/share/doc/git/contrib/diff-highlight/diff-highlight | less"
  git config --global pager.diff "perl /usr/share/doc/git/contrib/diff-highlight/diff-highlight | less"
  git config --global interactive.diffFilter "perl /usr/share/doc/git/contrib/diff-highlight/diff-highlight"

  # Colors
  git config --global color.ui true

  git config --global color.diff-highlight.oldNormal "red bold"
  git config --global color.diff-highlight.oldHighlight "red bold 52"
  git config --global color.diff-highlight.newNormal "green bold"
  git config --global color.diff-highlight.newHighlight "green bold 22"

  git config --global color.diff.meta "11"
  git config --global color.diff.frag "magenta bold"
  git config --global color.diff.commit "yellow bold"
  git config --global color.diff.old "red bold"
  git config --global color.diff.new "green bold"
  git config --global color.diff.whitespace "red reverse"
fi

