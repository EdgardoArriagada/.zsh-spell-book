reactCreator() {
  local TEMPLATES_PATH=${ZSB_DIR}/src/spells/react/react-creator/templates
  local templateName="$1"

  local newComponentPath="$2"
  local newComponentName="$(basename $newComponentPath)"

  local inputFlags="$3"
  local GROUP_FLAGS='o'

  ! ${zsb}_areFlagsInGroup "$inputFlags" "$GROUP_FLAGS" && return 1

  if [ -z "$newComponentPath" ]; then
    echo "${ZSB_ERROR} You must provide a React Component directory path"
    return 1
  fi


  if [ -d "$newComponentPath" ] && [[ "$inputFlags" != *'o'* ]] ; then
    echo "${ZSB_ERROR} Can't overwrite already existing component, use $(hl "-o") flag to do it anyway"
    return 1
  fi


  local replaceRegex="s/TemplateComponent/${newComponentName}/g"
  local chosenTemplate=${TEMPLATES_PATH}/${templateName}

  # Create components
  mkdir --parents "$newComponentPath" &&
    sed "$replaceRegex" ${chosenTemplate}/TemplateComponent.tsx \
    > ${newComponentPath}/${newComponentName}.tsx &&
    sed "$replaceRegex" ${chosenTemplate}/index.ts \
    > ${newComponentPath}/index.ts


  if [ $? ]; then
    echo "${ZSB_SUCCESS} New component -> ${newComponentPath}/$(hl "${newComponentName}.tsx")"
  fi
}

alias createreactcomponent="reactCreator FC"
alias createnextpage="reactCreator NEXTPAGE"
