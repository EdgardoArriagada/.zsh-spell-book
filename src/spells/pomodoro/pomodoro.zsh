# send instructions to pomodoro tmux session without attaching to it
pomodoro() (
  local this=$0
  local inputTime=$1
  shift 1
  local cbMessage=${*}

  ${this}.assertPomodoroNotRunning() {
    if `pdoro --is-counter-running`
      then ${zsb}.throw 'Pomodoro already running.'
    fi
  }

  ${this}.assertInputTime() {
    if ! `pdoro --is-valid-time ${inputTime}`
      then ${zsb}.throw 'Invalid input time.'
    fi
  }


  ${this}.decorateCbMessage() {
    if [[ -z "$cbMessage" ]]
      then cbMessage="Session for ${inputTime}"
      else cbMessage="${cbMessage} for ${inputTime}"
    fi
  }

  ${this}.beginPomodoro() {
    pdoro_cb --phase started $cbMessage

    # redirect only stdout, keep stderr for error handling
    pdoro -t $inputTime -c "pdoro_cb --phase ended ${cbMessage}" 1> /dev/null
  }

  { # main
    ${this}.assertInputTime

    ${this}.assertPomodoroNotRunning


    ${this}.decorateCbMessage
    ${this}.beginPomodoro

    # ${zsb}.success 'Pomodoro started.'
  }
)
