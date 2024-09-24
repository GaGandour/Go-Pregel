#/bin/bash

# Get arguments

DEBUG=false
FAILURE_STEP=-1
CHECKPOINT_FREQUENCY=-1
TEST=false

for arg in "$@"
do
  case $arg in
    -h|--help)
    echo "Usage: ./start_docker.sh -num_workers=<number of workers> -graph_file=<graph input file>"
    echo "The graph file is the relative path from the graphs folder."
    echo "Optional arguments:"
    echo "  -debug: Run in debug mode. This makes the pregel program to register the graph state in every superstep."
    echo "  -failure_step=<step number>: Simulate a failure in one of the workers at the specified step number. The worker will not be restarted and the computation will continue. The step number should be a positive integer."
    echo "  -test: Run the program in test mode. This will not open the graph visualization.\n"
    echo "Example 1: ./start_docker.sh -num_workers=3 -graph_file=graph1.json"
    echo "Example 2: ./start_docker.sh -num_workers=3 -graph_file=graph1.json -failure_step=5"
    echo "Example 3: ./start_docker.sh -num_workers=3 -graph_file=graph1.json -debug"
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
  case $arg in
    -failure_step=*)
    FAILURE_STEP="${arg#*=}"
    shift
    ;;
  esac
  case $arg in
    -checkpoint_frequency=*)
    CHECKPOINT_FREQUENCY="${arg#*=}"
    shift
    ;;
  esac
  case $arg in
    -test=*)
    TEST="${arg#*=}"
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

sh build_docker_image.sh

echo "Cleaning outputs from other pregel runs..."
./clean_outputs.sh # Clean previous outputs
echo "Finished cleaning outputs from other pregel runs"

cd ..

if [ "$DEBUG" = true ]; then
  echo "Running in debug mode"
  python3 python-scripts/write_docker_compose.py \
      --num_workers=$NUM_WORKERS \
      --graph_file=$GRAPH_FILE \
      --failure_step=$FAILURE_STEP \
      --checkpoint_frequency=$CHECKPOINT_FREQUENCY \
      --debug \
      > ../docker-compose.yml
else
  python3 python-scripts/write_docker_compose.py \
      --num_workers=$NUM_WORKERS \
      --graph_file=$GRAPH_FILE \
      --failure_step=$FAILURE_STEP \
      --checkpoint_frequency=$CHECKPOINT_FREQUENCY \
      > ../docker-compose.yml
fi

cd ..

# Create the folder structure for the output graph
graph_output_folder_structure=$(dirname "$GRAPH_FILE")
mkdir -p "src/output_graphs/$graph_output_folder_structure"
echo "Created folder structure: src/output_graphs/$graph_output_folder_structure"

docker-compose -f docker-compose.yml up -d
echo "Starting Pregel with $NUM_WORKERS workers on file $GRAPH_FILE"
docker attach pregel-master
echo "Stopping Pregel containers"
cd scripts/execution
sh ./stop_docker_containers.sh

if [ "$TEST" = false ]; then
    cd ../..
    source venv/bin/activate
    cd visualization
    python3 draw_graph.py --output_file=$GRAPH_FILE
    deactivate
    open graph.html
    cd ..
fi
