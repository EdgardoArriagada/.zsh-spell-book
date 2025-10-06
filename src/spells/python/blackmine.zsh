blackmine() {
    pea
    local gitFiles=( $(${zsb}.getGitFiles) )
    local validFiles=()

    for file in "${gitFiles[@]}"; do
        if [[ -f $file && $file == *.py ]]; then
            validFiles+=("$file")
        fi
    done

    black "${validFiles[@]}"
}
