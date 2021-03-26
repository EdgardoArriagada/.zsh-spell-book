# Create commands with composed arguments like
# mycommand -xyz
#
# paste the following code (adjusting the possition of flags var)
#
# {
#     local flags=$1 # or the possition where flags args goes
#     local GROUP_FLAGS='xyx'
#
#     ! ${zsb}.areFlagsInGroup "$flags" "$GROUP_FLAGS" && return 1
# }
#
# and ask if flags contains a char like:
# if [[ "$flags" == *'x'* ]]; then {...}
# if [[ "$flags" == *'y'* ]]; then {...}
# if [[ "$flags" == *'z'* ]]; then {...}
#
# or write this helper function to improve readability
# inputFlagsContains() [[ "$inputFlags" == *"$1"* ]]
#
# then you can call your command like this:
# mycommand -xyz
# mycommand -zyx
# mycommand -x
# mycommand -xz
# etc.

${zsb}.areFlagsInGroup() (
  local this="$0"
  local inputFlags="$1"
  local groupFlags="$2"

  ${this}.main() {
    if ! ${this}.areFlagsValid; then
      return 0
    fi

    if ! ${this}.doesInputFlagsBelongsToGroupFlags; then
      ${this}.throwUnknowFlagsException; return $?
    fi
  }

  ${this}.areFlagsValid() {
    [ -z "$inputFlags" ] && return 1

    local FLAG_REGEX="^-[a-z]+$"

    ${zsb}.doesMatch "$inputFlags" "$FLAG_REGEX"
  }

  ${this}.doesInputFlagsBelongsToGroupFlags() {
    local groupFlagsRegex="^-[${groupFlags}]+$"
    ${zsb}.doesMatch "$inputFlags" "$groupFlagsRegex"
  }

  ${this}.throwUnknowFlagsException() {
    echo "${ZSB_ERROR} One or more unknown flags in the list"
    return 1
  }

  ${this}.main "$@"
)

