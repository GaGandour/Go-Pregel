#/bin/bash

# Get arguments

VERBOSE=false

for arg in "$@"
do
    case $arg in
        -h|--help)
        echo "Usage: ./compare_results.sh -file_path=<algorithm name>/<graph name>"
        echo "Optional arguments:"
        echo "  -verbose: A little more explanation about the errors is given.\n"
        echo "Example 1: ./compare_results.sh -file_path=topological_sort/graph1.json"
        echo "Example 2: ./compare_results.sh -file_path=topological_sort/graph1.json -verbose"
        exit 0
        ;;
    esac
    case $arg in
        -algorithm=*)
        ALGORITHM="${arg#*=}"
        shift
        ;;
    esac
    case $arg in
        -verbose)
        VERBOSE=true
        shift
        ;;
    esac
done

if [ -z "$ALGORITHM" ]
  then
    echo "Missing algorithm. Run ./test_pregel_algorithm.sh with -h or --help for more information on the necessary arguments."
    exit 1
fi

cd ../python-scripts
source ../../venv/bin/activate
if [ $VERBOSE = true ]
then
    python3 graph_comparison.py --algorithm=$ALGORITHM --verbose
else
    python3 graph_comparison.py --algorithm=$ALGORITHM
fi
deactivate
cd ../auxiliary/

