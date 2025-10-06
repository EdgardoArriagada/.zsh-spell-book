blackmine() {
    pea
    local branch=""
    local gitFiles=()
    local validFiles=()

    # Parse options using zparseopts
    local -A opts
    zparseopts -D -A opts -- b: -branch:

    # Extract branch value from either -b or --branch
    branch=${opts[-b]:-${opts[--branch]}}

    # Get git files based on whether branch is provided
    if [[ -n "$branch" ]]; then
        gitFiles=( $(git diff --name-only "$branch" 2>/dev/null) )
        if [[ $? -ne 0 ]]; then
            echo "Error: Branch '$branch' not found or git diff failed"
            return 1
        fi
    else
        gitFiles=( $(${zsb}.getGitFiles) )
    fi

    # Filter for existing python files
    for file in "${gitFiles[@]}"; do
        if [[ -f $file && $file == *.py ]]; then
            validFiles+=("$file")
        fi
    done

    black "${validFiles[@]}"
}

# Completion function for blackmine
_blackmine() {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    # If this is the first word being completed, offer -b option
    if [[ $CURRENT -eq 2 && -z $words[2] ]]; then
        _describe 'options' '(-b)'
        return 0
    fi

    _arguments -C \
        '(-b --branch)'{-b,--branch}'[specify branch to compare against]:branch:_git_branch_names' \
        '*:file:_files'

    return 0
}

# Custom completion for git branch names
_git_branch_names() {
    local branches
    branches=($(git branch --format='%(refname:short)' 2>/dev/null))
    _describe 'branch' branches
}

compdef _blackmine blackmine
