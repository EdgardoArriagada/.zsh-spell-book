# example: $ echo "Text that is $(hl "highlighted") shows in color"
hl() {
  echo -n "${ZSB_SHL}"
  echo -E -n - "$@"
  echo -n "${ZSB_EHL}"
}
