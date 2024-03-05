${zsb}.getNoVerifyProp() {
  if [[ -f .pre-commit-config.yaml || -f .huskyrc || -d .husky ]]
    then <<< '--no-verify:Skip pre-commit verifications'
  fi
}
