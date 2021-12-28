qa() (
  local file=${1:?'You must provide a file'}
  local resultFile="${file}.result"

  if [[ -f $resultFile ]]; then
    ${zsb}.prompt "Do you want to delete $(hl $resultFile)? [Y/n]"
    ${zsb}.confirmMenu.continued && rm $resultFile
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

  local positions=($(seq 1 ${numQuestions}))
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
        questions[$currentPosition]="${item}"
        continue
      fi
      questions[$currentPosition]="${questions[$currentPosition]}\n${item}"
      continue
    fi

    if [[ "$state" = "$A" ]]; then
      if "$justSwitched"; then
        justSwitched=false
        answers[$currentPosition]="${item}"
        continue
      fi
      answers[$currentPosition]="${answers[$currentPosition]}\n${item}"
      continue
    fi
  done

  local -A badQuestions
  local -A badAnswers
  local indexes=()

  for i in {$numQuestions..1}; do
    printf "\n"
    printf "${questions[$i]}\n"
    printf "\n"
    ${zsb}.prompt "Press any key to continue"
    read -k1 -s
    printf "\n"
    printf "${answers[$i]}\n"
    printf "\n"
    ${zsb}.prompt "Was your answer correct? $(hl "[Y/n]")"
    while true; do
      read yn

      case $yn in
        [Yy]*) break ;;
        [Nn]*)
          badQuestions[$i]="${questions[$i]}"
          badAnswers[$i]="${answers[$i]}"
          indexes+=( $i )
          break ;;
        *) ${zsb}.prompt "Please answer yes or no" ;;
      esac
    done
  done

  for i in {$#badQuestions..1}; do
    local index=${indexes[$i]}
    echo "$Q" >> $resultFile
    echo "${badQuestions[$index]}" >> $resultFile
    echo "$A" >> $resultFile
    echo "${badAnswers[$index]}" >> $resultFile
    echo " " >> $resultFile
  done
)
