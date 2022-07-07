hasRust() {
  print " "; printAndRun "rustc --version"
  print " "; printAndRun "cargo --version"
  print " "; printAndRun "rustup --version"
  print " "
  ${zsb}.info "run `hl 'rustup update stable'` to update to the latest stable release"
}
