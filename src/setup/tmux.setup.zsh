${zsb}.safeLink tmux/tmux.conf ~/.tmux.conf

# Put tmux at the beggining of ~/.zshrc 
if (( $ZSB_MACOS ))
  then sed -i '.bak' '1s/^/[[ -z "$TMUX" ]] \&\& tmux new -A -s main\'$'\n/g' ~/.zshrc
  else sed -i '1i[[ -z "$TMUX" ]] \&\& tmux new -A -s main' ~/.zshrc
fi




