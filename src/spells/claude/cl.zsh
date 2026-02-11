cl() { claude --dangerously-skip-permissions $@; }

alias clp='cl --model opus --permission-mode plan'

hisIgnore cl clp
