# Create commands with composed arguments like
# mycommand -xyz
#
# paste the following code (adjusting the possition of flags var)
#
# {
#     local flags=$1 # or the possition where flags args goes
#     local GROUP_FLAGS='xyx'
#
#     ! areFlagsInGroup "$flags" "$GROUP_FLAGS" && return 1
# }
#
# and ask if flags contains a char like:
# if [[ "$flags" == *'x'* ]]; then {...}
# if [[ "$flags" == *'y'* ]]; then {...}
# if [[ "$flags" == *'z'* ]]; then {...}
#
# or write this helper function to improve readability
# inputFlagsContains() return $(([[ "$inputFlags" == *"$1"* ]]))
#
# then you can call your command like this:
# mycommand -xyz
# mycommand -zyx
# mycommand -x
# mycommand -xz
# etc.

areFlagsInGroup() {
  local flags="$1"
  local groupFlags="$2"

  if [ -z "$groupFlags" ]; then
    echo "${ZSB_ERROR} You must provide at least one possible flag as second arg. IE $(hl "xyz")"
    return 1
  fi

  [ -z "$flags" ] && return 0

  local flagsRegex="^-[${groupFlags}]+$"
  if ! doesMatch "$flags" "$flagsRegex"; then
    echo "${ZSB_ERROR} One or more unknown flags in the list"
    return 1
  fi

  return 0
}
