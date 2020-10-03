# Docker
alias dk='docker'
alias dkh='echorun docker --help'

# Docker Container
alias dkc='docker container'
alias dkch='echorun docker container --help'
alias dkr='docker container run'
alias dkrm='docker container rm'
alias dkll='echorun docker container ls'
alias dkl='echorun docker container ls -a'
alias dkle="echorun docker container ls \
  -f 'status=exited' \
  -f 'status=created' \
  -f 'status=restarting' \
  -f 'status=removing' \
  -f 'status=paused' \
  -f 'status=dead' \
  "

# Docker Image(s)
alias dki='docker image'
alias dkih='echorun docker image --help'
alias dkis='docker images'
alias dkil='echorun docker image ls'
alias dkrmi='docker image rm'
alias dkirm='dkrmi'
alias dkrmiall='docker rmi $(docker images -a -q)'
alias dkirmall='dkrmiall'

# Docker Volume
alias dkv='docker volume'
alias dkvh="echorun docker volume --help"
alias dkvl='echorun docker volume ls'
alias dkvrm="docker volume rm"

# Docker Network
alias dkn='docker network'
alias dknh='echorun docker network --help'
alias dknl='echorun docker network ls'

# Docker Service
alias dks='docker service'
alias dksh='echorun docker service --help'
alias dksl='echorun docker service ls'

# Docker Compose
alias dco='docker-compose'
alias dcoh='echorun docker-compose --help'
alias dcob='docker-compose build'
alias dcol='echorun docker-compose ps'
alias dcols='echorun docker-compose ps --services'
alias dcorm='docker-compose rm'
alias dcor='docker-compose run'
alias dcok='docker-compose kill'
alias dcou='echorun docker-compose up'
alias dcoud='echorun docker-compose up -d'

# Docker Machine
alias dkm='docker-machine'
alias dkmh='echorun docker-machine --help'
alias dkml='echorun docker-machine ls'
