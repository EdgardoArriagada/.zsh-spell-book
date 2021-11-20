${zsb}.formatBranch() {

  # [Sign/Data] text-data
  local signData=$(echo $1 | cut -d_ -f1)
  local textData=${1#*_}

  local branchType=$(echo $signData | cut -d/ -f1)
  local branchTicket=${signData#*/}

  local formattedText=$(echo $textData | tr '-' ' ' | sed 's/\w/\u&/')

  echo "[ ${(C)branchType} ${branchTicket} ] ${formattedText}"
}

createPr() {
  local formattedBranch=$(${zsb}.formatBranch $(git branch --show-current))
  local defaultBranch=$(git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@')

  gh pr create \
    --draft \
    --title "$formattedBranch" \
    --body 'WIP' \
    --base $defaultBranch
}
