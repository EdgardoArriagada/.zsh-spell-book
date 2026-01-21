grh() { git reset -q -- "$@" && ${zsb}.gitStatus; }

hisIgnore grh

compdef "_${zsb}.gitUnrepeat 'staged'" grh

