if [ -z "$1" ]
  then
    echo "No argument supplied."
    echo "Usage: ./start_docker.sh <number of workers> <graph input file>"
    exit 1
fi
if [ -z "$2" ]
  then
    echo "Only one argument supplied."
    echo "Usage: ./start_docker.sh <number of workers> <graph input file>"
    exit 1
fi

python3 write_docker_compose.py $1 $2 > ../docker-compose.yml
cd ..
mkdir -p output_graphs
docker-compose -f docker-compose.yml up -d
echo "Starting Pregel with $1 workers on file $2"
docker attach pregel-master
echo "Stopping Pregel containers"
cd scripts
sh ./stop_docker.sh
cd ..
cd visualization
python3 draw_graph.py ../src/output_graphs/output_graph.json