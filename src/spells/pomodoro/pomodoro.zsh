# send instructions to pomodoro tmux session without attaching to it
pomodoro() (
  local this=$0
  local inputTime=$1
  shift 1
  local pomodoroLabel=${*}

  local cbMessage

  ${this}.validatePomodoroNotRunning() {
    if ${this}.isPomodoroRunning; then
      ${zsb}.throw 'A pomodoro timer is already running.'
    fi
  }


  ${this}.setBeginString() {
    if [[ -z "$label" ]]
      then cbMessage="Session for ${inputTime}"
      else cbMessage="${label} for ${inputTime}"
    fi
  }

  ${this}.beginPomodoro() {
    pdoro_cb --phase started ${cbMessage}

    # redirect only stdout, keep stderr for error handling
    pdoro -t $inputTime -c "pdoro_cb --phase ended ${cbMessage}" 1> /dev/null
  }

  { # main
    ## TODO: validate input time

    ## TODO: validate pomodoro not running


    ${this}.setBeginString
    ${this}.beginPomodoro

    # ${zsb}.success 'Pomodoro started.'
  }
)
