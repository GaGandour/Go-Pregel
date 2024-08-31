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
        -file_path=*)
        FILE_PATH="${arg#*=}"
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

if [ -z "$FILE_PATH" ]
  then
    echo "Missing file_path. Run ./test_pregel_algorithm.sh with -h or --help for more information on the necessary arguments."
    exit 1
fi

cd ../python-scripts
source ../../venv/bin/activate
if [ $VERBOSE = true ]
then
    echo 'verbose'
    python3 graph_comparison.py --graph_file=$FILE_PATH --verbose
else
    echo 'not verbose'
    python3 graph_comparison.py --graph_file=$FILE_PATH
fi
deactivate
cd ../execution

