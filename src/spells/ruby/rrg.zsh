function rrg() {
  # Fix pipeline preventing
  # from loading rails
  # On current terminal session
  eval "rails -v >/dev/null 2>&1"

  eval "rails routes | grep $@"

  alias rrg="rails routes | grep"
  unfunction rrg
}

