docker kill $(docker ps -q)
docker rm -f $(docker ps -aq)
docker rmi $(docker images "dev-*" -q)
