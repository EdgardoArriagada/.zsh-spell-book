tmheaderbar() {
  local winId
  if [[ $1 =~ ^@[0-9]+$ ]]
    then winId=$1; shift 1
    else winId=`tmux display -p "#{window_id}"`
  fi

  tmux split-window -vb -p 6 -t ${winId}.0 -P -F "#{pane_id}" "print_centered \"${*}\" | less" \; \
    select-pane -t ${winId}.1
}
