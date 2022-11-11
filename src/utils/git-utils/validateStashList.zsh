${zsb}.isEmptyStashList() [[ -z "`git stash list`" ]]
${zsb}.validateStashList() { ${zsb}.isEmptyStashList && ${zsb}.cancel 'Stash list is empty.'; }
