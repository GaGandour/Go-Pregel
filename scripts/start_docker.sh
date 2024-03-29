if [ -z "$1" ]
  then
    echo "No argument supplied."
    echo "Usage: ./start_docker.sh <number of workers>"
    exit 1
fi
# TODO: BUILD DOCKER-COMPOSE FILE
cd ..
docker-compose -f docker-compose.yml up -d
echo "Starting DII with $1 workers"
docker attach dii-master
echo "Stopping DII containers"
cd scripts
./stop_docker.sh