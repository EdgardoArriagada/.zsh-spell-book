improve-pr-description(){
  eval "gh pr view --json number --jq '.number'"
  zsb_improve-pr-description
}
