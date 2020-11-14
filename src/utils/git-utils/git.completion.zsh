# All completions must work from current working dir (not only from project root)
# Example usage: complete -C $ZSB_GIT_STAGED_GILES mycmd

ZSB_GIT_STAGED_FILES="${zsb}_isGitRepo && git status --short | grep '^[MARCD]' | sed s/^...//"
ZSB_GIT_UNSTAGED_FILES="${zsb}_isGitRepo && git status --short | grep -E '(^ [MARCD])|(^[MARCD]{2})' | sed s/^...//"
ZSB_GIT_UNTRACKED_FILES="${zsb}_isGitRepo && git status --short | grep -E '(^ \?)|(^\?{2})' | sed s/^...//"
ZSB_GIT_UNSTAGED_AND_UNTRACKED_FILES="${zsb}_isGitRepo && git status --short | grep -E '(^ [MARCD\?])|(^[MARCD\?]{2})' | sed s/^...//"
ZSB_GIT_MODIFIED_FILES="${zsb}_isGitRepo && git status --short | sed s/^...// || ls"
