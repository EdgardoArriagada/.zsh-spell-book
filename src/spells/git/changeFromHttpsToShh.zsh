changeFromHttpsToShh() {
  local username=`git config user.name`
  local reponame=`get_repo_name`

  if [[ -z "$reponame" || -z "$username" ]]; then
    ${zsb}.throw "Could not fetch reponame or username"
  fi

  git remote set-url origin git@github.com:$username/$reponame
}
