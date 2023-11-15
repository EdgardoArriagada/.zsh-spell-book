cdw() cds $ZSB_PROJECTS_DIR/$1
cdw2() cds ${ZSB_PROJECTS_DIR}2/$1
cdw3() cds ${ZSB_PROJECTS_DIR}3/$1
cdw4() cds ${ZSB_PROJECTS_DIR}4/$1

tmw() { tm $1 -c $ZSB_PROJECTS_DIR/$1 }
tmw2() { tm $1 -c ${ZSB_PROJECTS_DIR}2/$1 }
tmw3() { tm $1 -c ${ZSB_PROJECTS_DIR}3/$1 }
tmw4() { tm $1 -c ${ZSB_PROJECTS_DIR}4/$1 }

compdef "_${zsb}.singleCompC 'ls $ZSB_PROJECTS_DIR'" cdw tmw
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}2'" cdw2 tmw2
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}3'" cdw3 tmw3
compdef "_${zsb}.singleCompC 'ls ${ZSB_PROJECTS_DIR}4'" cdw4 tmw4
