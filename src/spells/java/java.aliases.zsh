alias configurejava='sudo update-alternatives --config java'
alias javaVersion="java -version 2>&1 | awk -F '\"' 'FNR==2 {print \$2}'"
