qa() (
  local file=${1:?Error: You must provide a file.}

  ${zsb}.doesMatch "$file" "\.qa$" || ${zsb}.throw "File has to be a $(hl ".qa")"

  local newFileName=$(echo "$file" | cut -d'.' -f'1' | tr -d -c a-zA-Z )

  integer fileVersion=$(echo "$file" | tr -d -c 0-9)
  [[ -z "fileVersion" ]] && fileVersion=0
  (( fileVersion++ ))

  local resultFile="${newFileName}_${fileVersion}.qa"

  if [[ -f $resultFile ]]; then
    ${zsb}.prompt "Do you want to delete $(hl $resultFile)? [Y/n]"
    ${zsb}.confirmMenu && rm $resultFile
  fi

  zmodload zsh/mapfile
  local fileLines=( "${(f)mapfile[$file]}" )
  local -A questions
  local -A answers

  local state # Q-- | A--
  local Q='Q--'
  local A='A--'
  local justSwitched=false

  integer numQuestions=0
  # Count Questions
  for item in $fileLines; do
    [[ "$item" = "$Q" ]] && (( numQuestions++ ))
  done

  local positions=($(seq 1 $numQuestions))
  positions=( $(shuf -e "${positions[@]}") )

  local currentPosition
  integer i=0
  for item in $fileLines; do
    if [[ "$item" = "$Q" ]]; then
      state="$Q"
      justSwitched=true
      (( i++ ))
      currentPosition="${positions[$i]}"
      continue
    fi

    if [[ "$item" = "$A" ]]; then
      state="$A"
      justSwitched=true
      continue
    fi

    if [[ "$state" = "$Q" ]]; then
      if "$justSwitched"; then
        justSwitched=false
        questions[$currentPosition]=$item
        continue
      fi
      questions[$currentPosition]="${questions[$currentPosition]}\n${item}"
      continue
    fi

    if [[ "$state" = "$A" ]]; then
      if "$justSwitched"; then
        justSwitched=false
        answers[$currentPosition]=$item
        continue
      fi
      answers[$currentPosition]="${answers[$currentPosition]}\n${item}"
      continue
    fi
  done

  local -A failedQuestions
  local -A failedAnswers
  local indexes=()

  for i in {1..$numQuestions}; do
    color 246 `hr`
    printCentered "[${i} / ${numQuestions}]"
    print " "
    printf "${questions[$i]}\n"
    read -k1 -s
    color 76 `hr`
    printf "${answers[$i]}\n"
    color 76 `hr`
    ${zsb}.prompt "Was your answer correct? $(hl "[Y/n]")"
    while true; do
      read yn

      case $yn in
        [Yy]*) break ;;
        [Nn]*)
          failedQuestions[$i]="${questions[$i]}"
          failedAnswers[$i]="${answers[$i]}"
          indexes+=( $i )
          break ;;
        *) ${zsb}.prompt "Please answer yes or no" ;;
      esac
    done
  done

  [[ "$#failedQuestions" = "0" ]] &&
    ${zsb}.success "You answered all the Questionary correctly!" && return 0

  for i in {1..$#failedQuestions}; do
    local index=${indexes[$i]}
    echo "$Q" >> $resultFile
    echo "${failedQuestions[$index]}" >> $resultFile
    echo "$A" >> $resultFile
    echo "${failedAnswers[$index]}" >> $resultFile
    echo " " >> $resultFile
  done

  ${zsb}.info "$(hl "./$resultFile") generated with $(hl $#failedQuestions) failing question(s)"
)
