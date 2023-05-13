hasRust() {
  <<< ''; printAndRun "rustc --version"
  <<< ''; printAndRun "cargo --version"
  <<< ''; printAndRun "rustup --version"
  <<< ''
  ${zsb}.info "run `hl 'rustup update stable'` to update to the latest stable release"
}
