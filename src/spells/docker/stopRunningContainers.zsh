stopRunningContainers() { 
  local containers=$(docker container ls -q) 
  if [ ! -z "$containers" ]; then 
    echo "${ZSB_INFO} Stopping all running containers..." 

    docker stop $(echo "$containers") && \
      echo "${ZSB_SUCCESS} All containers did halt.\n" && \
      return 0
  fi 

  echo "${ZSB_INFO} There are no running containers."
} 
