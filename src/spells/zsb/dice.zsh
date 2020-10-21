dice() (
  local sides="$1"
  local DEFAULT_SIDES=6

  main() {
    setDefaultSides

    if ! areSidesValid; then
      echo "${ZSB_ERROR} Invalid number of sides. $(it "{sides : sides ∈ Z and 2 ≤ sides ≤ ∞}")"
      return 1
    fi

    rollTheDice
  }

  setDefaultSides() {
    [ -z "$sides" ] && sides="$DEFAULT_SIDES"
  }

  areSidesValid() {
    return $(${zsb}_isInteger "$sides" && [ "$sides" -gt "1" ])
  }

  rollTheDice() {
    local result=$(generateRandom)
    echo "[${sides}]: $(hl "$result")"
  }

  generateRandom() {
    # changing seed is mandatory as base seed doesn't change in a subshell
    local randomNumber=$(head -1 /dev/urandom | od -N 1 | awk '{ print $2 }')
    echo $((randomNumber % ${sides} + 1))
  }

  main "$@"
)

complete -W "2 4 6 8 12" dice

alias die=dice
