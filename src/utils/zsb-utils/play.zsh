# Play comes with 'sox' in macos
${zsb}.play() {
  (( $ZSB_MACOS )) && \
    play $1 >/dev/null 2>&1 || \
    aplay $1 >/dev/null 2>&1
}
