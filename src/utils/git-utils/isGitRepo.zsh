${zsb}.isGitRepo() { [[ -d .git ]] || `git rev-parse --is-inside-work-tree 2>/dev/null`; }

