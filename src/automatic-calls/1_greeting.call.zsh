# surrounding list in parentheses creates a subshell [for a clear namespace]

(
  # to easily generate an ascii art, google for pages that convert text or images to ascii

  CL=$(((RANDOM % 6) + 31))
  SD=$((RANDOM % 2))
  COLOR="\033[${SD};${CL}m"
  NOCOL='\033[0m'

  echo " "
  echo "${COLOR}      ╦ ╦╔═╗╦  ╔═╗╔═╗╔╦╗╔═╗"
  echo "      ║║║║╣ ║  ║  ║ ║║║║║╣"
  echo "${NOCOL}      ╚╩╝╚═╝╩═╝╚═╝╚═╝╩ ╩╚═╝"
  echo " "
)
