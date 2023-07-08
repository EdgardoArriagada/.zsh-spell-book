dice() (
  local this=$0
  local sides=${1:=6}

  ${this}.validateSides() {
    ${this}.areSidesValid && return 0
    ${zsb}.throw "Invalid number of sides. $(it "{sides : sides ∈ Z and 2 ≤ sides ≤ ∞}")"
  }

  ${this}.areSidesValid() {
    ${zsb}.isInteger $sides && (( $sides > 1 ))
  }

  ${this}.rollTheDice() {
    local result=`${this}.generateRandom`
    <<< "[${sides}]: $(hl "$result")"
  }

  ${this}.generateRandom() {
    # changing seed is mandatory as base seed doesn't change in a subshell
    local randomNumber=`head -1 /dev/urandom | od -An | awk 'FNR == 1 { print $2 }'`
    <<< $((randomNumber % ${sides} + 1))
  }

  { # main
    ${this}.validateSides

    ${this}.rollTheDice
  }
)

complete -W "2 4 6 8 12" dice

alias die=dice
