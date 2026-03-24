hasRust() {
  if ! printAndRun "rustc --version"; then
    ${zsb}.info "Rust not found. Redirecting to `hl 'https://rustup.rs'`"
    loading 2
    zsb_open "https://rustup.rs"
    return 1
  fi

  printAndRun "cargo --version"
  printAndRun "rustup --version"
  ${zsb}.info "run `hl 'rustup update'` to update all toolchains and rustup itself"
}
