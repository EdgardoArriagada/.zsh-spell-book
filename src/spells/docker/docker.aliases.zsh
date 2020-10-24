# Docker
alias dk='docker'
alias dkh="printAndRun 'docker --help'"

# Docker Container
alias dkc='docker container'
alias dkch="printAndRun 'docker container --help'"
alias dkr='docker container run'
alias dkrm='docker container rm'
alias dkll="printAndRun 'docker container ls'"
alias dkl="printAndRun 'docker container ls -a'"
alias dkle="printAndRun \"docker container ls\
 -f 'status=exited'\
 -f 'status=created'\
 -f 'status=restarting'\
 -f 'status=removing'\
 -f 'status=paused'\
 -f 'status=dead'\"
"

# Docker Image(s)
alias dki='docker image'
alias dkih="printAndRun 'docker image --help'"
alias dkis='docker images'
alias dkil="printAndRun 'docker image ls'"
alias dkrmi='docker image rm'
alias dkirm='dkrmi'
alias dkrmiall='docker rmi $(docker images -a -q)'
alias dkirmall='dkrmiall'

# Docker Volume
alias dkv='docker volume'
alias dkvh="printAndRun 'docker volume --help'"
alias dkvl="printAndRun 'docker volume ls'"
alias dkvrm="docker volume rm"

# Docker Network
alias dkn='docker network'
alias dknh="printAndRun 'docker network --help'"
alias dknl="printAndRun 'docker network ls'"

# Docker Service
alias dks='docker service'
alias dksh="printAndRun 'docker service --help'"
alias dksl="printAndRun 'docker service ls'"

# Docker Compose
alias dco='docker-compose'
alias dcoh="printAndRun 'docker-compose --help'"
alias dcob='docker-compose build'
alias dcol="printAndRun 'docker-compose ps'"
alias dcols="printAndRun 'docker-compose ps --services'"
alias dcorm='docker-compose rm'
alias dcor='docker-compose run'
alias dcok='docker-compose kill'
alias dcou="printAndRun 'docker-compose up'"
alias dcoud="printAndRun 'docker-compose up -d'"

# Docker Machine
alias dkm='docker-machine'
alias dkmh="printAndRun 'docker-machine --help'"
alias dkml="printAndRun 'docker-machine ls'"
