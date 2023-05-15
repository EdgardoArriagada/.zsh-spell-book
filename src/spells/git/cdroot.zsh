cdroot() {
  local pathToDotGit=`git rev-parse --git-dir 2>/dev/null`
  [[ -n "$pathToDotGit" ]] && cds `dirname ${pathToDotGit}`
}

hisIgnore cdroot

_${zsb}.nocompletion cdroot

