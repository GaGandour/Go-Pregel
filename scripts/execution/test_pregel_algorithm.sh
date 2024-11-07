#/bin/bash

# Get arguments

FAILURE_STEP=-1
VERBOSE=false
CHECKPOINT_FREQUENCY=-1
SKIP_PREGEL=false

for arg in "$@"
do
    case $arg in
        -h|--help)
        echo "Usage: ./test_pregel_algorithm.sh -num_workers=<number of workers> -algorithm=<algorithm name>"
        echo "Optional arguments:"
        echo "  -failure_step=<step number>: Simulate a failure in one of the workers at the specified step number. The worker will not be restarted and the computation will continue. The step number should be a positive integer."
        echo "  -checkpoint_frequency=<frequency>: Checkpoint the graph state every <frequency> supersteps. The frequency should be a positive integer."
        echo "  -verbose: A little more explanation about the errors is given.\n"
        echo "  -skip_pregel: Skip the pregel execution and only compare the output files.\n"
        echo "Example 1: ./test_pregel_algorithm.sh -num_workers=3 -algorithm=topological_sort"
        echo "Example 2: ./test_pregel_algorithm.sh -num_workers=3 -algorithm=topological_sort -failure_step=5 -checkpoint_frequency=2"
        echo "Example 3: ./test_pregel_algorithm.sh -num_workers=3 -algorithm=topological_sort -verbose"
        echo "Example 4: ./test_pregel_algorithm.sh -num_workers=3 -algorithm=topological_sort -skip_pregel"
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
        -algorithm=*)
        ALGORITHM="${arg#*=}"
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
        -verbose)
        VERBOSE=true
        shift
        ;;
    esac
    case $arg in
        -skip_pregel)
        SKIP_PREGEL=true
        shift
        ;;
    esac
done

if [ -z "$NUM_WORKERS" ]
  then
    echo "Missing Number of Workers. Run ./test_pregel_algorithm.sh with -h or --help for more information on the necessary arguments."
    exit 1
fi
if [ -z "$ALGORITHM" ]
  then
    echo "Missing algorithm. Run ./test_pregel_algorithm.sh with -h or --help for more information on the necessary arguments."
    exit 1
fi

# Iterate through the list of graph files if skip pregel is not set
if [ $SKIP_PREGEL = false ]
then
    for filename in ../../graphs/${ALGORITHM}/*.json; do
        # Extract the graph name from the file name
        graph_name=$(basename $filename)
        graph_path="${ALGORITHM}/${graph_name}"
        # Execute pregel and wait for output files
        ./start_pregel.sh \
            -num_workers=$NUM_WORKERS \
            -graph_file=$graph_path \
            -failure_step=$FAILURE_STEP \
            -checkpoint_frequency=$CHECKPOINT_FREQUENCY \
            -test
    done
fi

cd ../auxiliary/
# Compare the output files
if [ $VERBOSE = true ]
then
    ./compare_results.sh -algorithm=$ALGORITHM -verbose
else
    ./compare_results.sh -algorithm=$ALGORITHM
fi
