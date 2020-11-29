dice() (
  local this="$0"
  local sides="$1"
  local DEFAULT_SIDES=6

  ${this}.main() {
    ${this}.setDefaultSides

    if ! ${this}.areSidesValid; then
      echo "${ZSB_ERROR} Invalid number of sides. $(it "{sides : sides ∈ Z and 2 ≤ sides ≤ ∞}")"
      return 1
    fi

    ${this}.rollTheDice
  }

  ${this}.setDefaultSides() {
    [ -z "$sides" ] && sides="$DEFAULT_SIDES"
  }

  ${this}.areSidesValid() {
    ${zsb}.isInteger "$sides" && [ "$sides" -gt "1" ]
  }

  ${this}.rollTheDice() {
    local result=$(${this}.generateRandom)
    echo "[${sides}]: $(hl "$result")"
  }

  ${this}.generateRandom() {
    # changing seed is mandatory as base seed doesn't change in a subshell
    local randomNumber=$(head -1 /dev/urandom | od -N 1 | awk '{ print $2 }')
    echo $((randomNumber % ${sides} + 1))
  }

  ${this}.main "$@"
)

complete -W "2 4 6 8 12" dice

alias die=dice
