sqlineOf() (
  local this="$0"

  local action="$1"
  local tableName="$2"
  local file="$3"

  ${this}.main() {
    if ! ${this}.areArgsValid "$@"; then
      ${this}.throwInvalidArgs; return $?
    fi

    case "$action" in
      'INSERT_INTO')
        ${this}.lineOfMatch "INSERT INTO \`${tableName}\`" ;;
      'CREATE_TABLE')
        ${this}.lineOfMatch "CREATE TABLE \`${tableName}\`" ;;
      *)
        echo "${ZSB_ERROR} Unhanlded action." ;;
    esac
  }

  ${this}.areArgsValid() [ "$#" -gt 2 ]

  ${this}.throwInvalidArgs() {
    echo "${ZSB_ERROR} You must probide an action, table name and a file name."
    return 1
  }

  ${this}.lineOfMatch() {
    grep -n "$1" ${file} | cut -f1 -d:
  }

  ${this}.main "$@"
)

_${zsb}.sqlineOf() {
  case $COMP_CWORD in
    1)
      COMPREPLY=( $(compgen -W "'INSERT_INTO' 'CREATE_TABLE'") )
      ;;
    2)
      COMPREPLY=( $(compgen -W "my_table") )
      ;;
    3)
      COMPREPLY=( $(compgen -f -X  "!*.sql") )
      ;;
  esac

}

complete -F _${zsb}.sqlineOf sqlineOf

