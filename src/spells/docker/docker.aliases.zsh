# Docker
alias dk='docker'
alias dkh="printAndRun 'docker --help'"

# Docker Container
alias dkc='docker container'
alias dkch="printAndRun 'docker container --help'"
alias dkr='docker container run'
alias dkb='docker build'
alias dkrm='docker container rm'
alias dkl="printAndRun 'docker container ls'"
alias dkll="printAndRun 'docker container ls -a'"
alias dkle="printAndRun \"docker container ls \
 -f 'status=exited' \
 -f 'status=created' \
 -f 'status=restarting' \
 -f 'status=removing' \
 -f 'status=paused' \
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

# Docker Volume
alias dkv='docker volume'
alias dkvh="printAndRun 'docker volume --help'"
alias dkvl="printAndRun 'docker volume ls'"
alias dkvrm="docker volume rm"
alias dkvrmall='docker volume rm `docker volume ls -q`'

# Docker Network
alias dkip="copythis '172.17.0.1'"
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
alias dcod="printAndRun 'docker-compose down'"
alias dcos="printAndRun 'docker-compose stop'"
alias dcorm='docker-compose rm'
alias dcok='docker-compose kill'

# Docker Machine
alias dkm='docker-machine'
alias dkmh="printAndRun 'docker-machine --help'"
alias dkml="printAndRun 'docker-machine ls'"
