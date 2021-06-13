reactCreator() {
  local -r TEMPLATES_PATH=${ZSB_DIR}/src/spells/react/react-creator/templates
  local -r templateName="$1"

  local -r newComponentPath=${2:?'You must provide a React Component directory path'}
  local -r newComponentKebab="$(basename $newComponentPath)"
  local -r newComponentUpperCamelCase="$(echo "$newComponentKebab" | sed -r 's/(^|-)([a-z])/\U\2/g')"

  local -r inputFlags="$3"
  local -r GROUP_FLAGS='o'

  ! ${zsb}.areFlagsInGroup "$inputFlags" "$GROUP_FLAGS" && return 1

  if [[ -d "$newComponentPath" ]] && [[ "$inputFlags" != *'o'* ]] ; then
    ${zsb}.throw "Can't overwrite already existing component, use $(hl "-o") flag to do it anyway"
  fi

  local -r chosenTemplate=${TEMPLATES_PATH}/${templateName}

  # Create components
  mkdir --parents "$newComponentPath" &&
    sed "s/TemplateComponent/${newComponentUpperCamelCase}/g" ${chosenTemplate}/TemplateComponent.tsx \
    > ${newComponentPath}/${newComponentKebab}.tsx &&
    sed "s/TemplateComponent/${newComponentKebab}/g" ${chosenTemplate}/index.tsx \
    > ${newComponentPath}/index.tsx


  if [[ "$?" == "0" ]]; then
    echo "${ZSB_SUCCESS} New component -> ${newComponentPath}/$(hl "${newComponentKebab}.tsx")"
  fi
}

alias createreactcomponent="reactCreator FC"
alias createnextpage="reactCreator NEXTPAGE"

