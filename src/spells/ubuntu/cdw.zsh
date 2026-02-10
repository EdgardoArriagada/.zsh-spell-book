${zsb}.cdw_generic() {
    local project_dir=$1
    local project=$2
    if [[ -z "$project" ]]; then
      project=$(ls $project_dir | rg -v '_gitworktree$' | fzf)
      [[ -n "$project" ]] && print -z "cdw $project"
    else
      cds $project_dir/$project
    fi
}

${zsb}.tmw_generic() {
    local project_dir=$1
    local tmnum=$2
    local project=$3
    if [[ -z "$project" ]]; then
      project=$(ls $project_dir | rg -v '_gitworktree$' | fzf)
      [[ -n "$project" ]] && print -z "tmw${tmnum} $project"
    else
      local suffix=""
      [[ -n "$tmnum" ]] && suffix="-${tmnum}"
      tm ${project}${suffix} -c $project_dir/$project
    fi
}

cdw() { ${zsb}.cdw_generic $ZSB_PROJECTS_DIR $1 }
cdw2() { ${zsb}.cdw_generic ${ZSB_PROJECTS_DIR}2 $1 }
cdw3() { ${zsb}.cdw_generic ${ZSB_PROJECTS_DIR}3 $1 }
cdw4() { ${zsb}.cdw_generic ${ZSB_PROJECTS_DIR}4 $1 }

tmw() { ${zsb}.tmw_generic $ZSB_PROJECTS_DIR "" $1 }
tmw2() { ${zsb}.tmw_generic ${ZSB_PROJECTS_DIR}2 2 $1 }
tmw3() { ${zsb}.tmw_generic ${ZSB_PROJECTS_DIR}3 3 $1 }
tmw4() { ${zsb}.tmw_generic ${ZSB_PROJECTS_DIR}4 4 $1 }

compdef "_${zsb}.singleCompC 'ls $ZSB_PROJECTS_DIR | rg -v _gitworktree\$'" cdw tmw
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}2 | rg -v _gitworktree\$'" cdw2 tmw2
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}3 | rg -v _gitworktree\$'" cdw3 tmw3
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}4 | rg -v _gitworktree\$'" cdw4 tmw4
