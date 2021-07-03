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

  local -r templateTestFile=${chosenTemplate}/TemplateComponent.test.tsx
  local -r resultTestFile=${newComponentPath}/${newComponentKebab}.test.tsx

  local -r templateComponent=${chosenTemplate}/TemplateComponent.tsx
  local -r resultComponent=${newComponentPath}/${newComponentKebab}.tsx

  local -r templateIndexFile=${chosenTemplate}/index.tsx
  local -r resultIndexFile=${newComponentPath}/index.tsx

  local -r upperToUpper=s/TemplateComponent/${newComponentUpperCamelCase}/g
  local -r kebabToKebab=s/template-component/${newComponentKebab}/g

  # Create components
  mkdir --parents "$newComponentPath" && \

    # Test File
    sed "$upperToUpper" $templateTestFile > $resultTestFile && \
    sed --in-place "$kebabToKebab" $resultTestFile && \

    # Main Component
    sed "$upperToUpper" $templateComponent > $resultComponent && \

    # Index File
    sed "$kebabToKebab" $templateIndexFile > $resultIndexFile


  if [[ "$?" == "0" ]]; then
    echo "${ZSB_SUCCESS} New component -> ${newComponentPath}/$(hl "${newComponentKebab}.tsx")"
  fi
}

alias createreactcomponent="reactCreator FC"
alias createnextpage="reactCreator NEXTPAGE"

