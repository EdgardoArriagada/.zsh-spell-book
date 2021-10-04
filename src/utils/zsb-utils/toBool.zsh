Bool() {
  local input="$1"

  if ${zsb}.isFloat "$input"; then
    [[ "$1" = "0" ]] && echo true && return
    echo false && return
  fi

  [[ -z "$input" ]] && echo false && return

  echo true
}

### TRUE
# Bool 0 # true
# Bool '0' # true
# Bool ' ' # true
# Bool " " # true
# Bool foo # true
# Bool "bar" # true

### FALSE
# Bool 1 # false
# Bool '1' # false
# Bool '' # false
# Bool "" # false
