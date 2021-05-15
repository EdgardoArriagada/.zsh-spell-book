# zsh completions does not read aliases, so tmls is a function
tmls() { tmux ls 2>&1 | cut -d':' -s -f1 }

