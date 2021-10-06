# 1=input, 2=regex
${zsb}.doesMatch() [[ "$1" =~ "$2" ]]

${zsb}.isInteger() [[ "$1" = <-> ]]

${zsb}.isFloat() [[ "$1" =~ "^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$" ]]

${zsb}.isUrl() [[ "$1" =~ "^http[s]?:\/{2}" ]]

