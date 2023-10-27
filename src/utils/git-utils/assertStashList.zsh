${zsb}.isEmptyStashList() [[ -z "`git stash list`" ]]
${zsb}.assertStashList() { ${zsb}.isEmptyStashList && ${zsb}.cancel 'Stash list is empty.'; }
