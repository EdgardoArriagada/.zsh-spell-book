${zsb}.cdw_generic() {
    local project_dir=$1
    if [[ -z "$2" ]]; then
        local dir=$(ls $project_dir | fzf)
        [[ -n "$dir" ]] && cds $project_dir/$dir
    else
        cds $project_dir/$2
    fi
}

${zsb}.tmw_generic() {
    local project_dir=$1
    local suffix=$2
    if [[ -z "$3" ]]; then
        local dir=$(ls $project_dir | fzf)
        [[ -n "$dir" ]] && tm ${dir}${suffix} -c $project_dir/$dir
    else
        tm ${3}${suffix} -c $project_dir/$3
    fi
}

cdw() { ${zsb}.cdw_generic $ZSB_PROJECTS_DIR $1 }
cdw2() { ${zsb}.cdw_generic ${ZSB_PROJECTS_DIR}2 $1 }
cdw3() { ${zsb}.cdw_generic ${ZSB_PROJECTS_DIR}3 $1 }
cdw4() { ${zsb}.cdw_generic ${ZSB_PROJECTS_DIR}4 $1 }

tmw() { ${zsb}.tmw_generic $ZSB_PROJECTS_DIR "" $1 }
tmw2() { ${zsb}.tmw_generic ${ZSB_PROJECTS_DIR}2 "-2" $1 }
tmw3() { ${zsb}.tmw_generic ${ZSB_PROJECTS_DIR}3 "-3" $1 }
tmw4() { ${zsb}.tmw_generic ${ZSB_PROJECTS_DIR}4 "-4" $1 }

compdef "_${zsb}.singleCompC 'ls $ZSB_PROJECTS_DIR'" cdw tmw
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}2'" cdw2 tmw2
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}3'" cdw3 tmw3
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}4'" cdw4 tmw4
