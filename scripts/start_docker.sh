#/bin/bash

# Get arguments

DEBUG=false

for arg in "$@"
do
  case $arg in
    -h|--help)
    echo "Usage: ./start_docker.sh -num_workers=<number of workers> -graph_file=<graph input file>"
    echo "Optional arguments:"
    echo "  -debug: Run in debug mode. This makes the pregel program to register the graph state in every superstep.\n"
    echo "Example 1: ./start_docker.sh -num_workers=3 -graph_file=graph1.json"
    echo "Example 1: ./start_docker.sh -num_workers=3 -graph_file=graph1.json -debug"
    exit 0
    ;;
  esac
  case $arg in
    -num_workers=*)
    NUM_WORKERS="${arg#*=}"
    shift
    ;;
  esac
  case $arg in
    -graph_file=*)
    GRAPH_FILE="${arg#*=}"
    shift
    ;;
  esac
  case $arg in
    -debug)
    DEBUG=true
    shift
    ;;
  esac
done

if [ -z "$NUM_WORKERS" ]
  then
    echo "Missing Number of Workers. Run ./start_docker with -h or --help for more information on the necessary arguments."
    exit 1
fi
if [ -z "$GRAPH_FILE" ]
  then
    echo "Missing Graph File. Run ./start_docker with -h or --help for more information on the necessary arguments."
    exit 1
fi

sh build_image.sh

if [ "$DEBUG" = true ]; then
  echo "Running in debug mode"
  python3 write_docker_compose.py $NUM_WORKERS $GRAPH_FILE -debug > ../docker-compose.yml
else
  python3 write_docker_compose.py $NUM_WORKERS $GRAPH_FILE > ../docker-compose.yml
fi

cd ..
mkdir -p src/output_graphs
docker-compose -f docker-compose.yml up -d
echo "Starting Pregel with $NUM_WORKERS workers on file $GRAPH_FILE"
docker attach pregel-master
echo "Stopping Pregel containers"
cd scripts
sh ./stop_docker.sh
cd ..
source venv/bin/activate
cd visualization
python3 draw_graph.py ../src/output_graphs/output_graph.json
deactivate
open graph.html
cd ..