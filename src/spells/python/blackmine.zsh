blackmine() {
    pea
    local gitFiles=( $(${zsb}.getGitFiles) )
    black "${gitFiles[@]}"
}
