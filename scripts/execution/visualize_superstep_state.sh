#/bin/bash

for arg in "$@"
do
    case $arg in
        -h|--help)
        echo "Usage: ./visualize_superstep_state.sh -superstep=<superstep number>"
        echo "Or: ./visualize_superstep_state.sh -output_file=<output file name>"
        exit 0
        ;;
    esac
    case $arg in
        -superstep=*)
        SUPERSTEP="${arg#*=}"
        shift
        ;;
    esac
    case $arg in
        -output_file=*)
        OUTPUT_FILE="${arg#*=}"
        shift
        ;;
    esac
done

cd ../..
source venv/bin/activate
cd visualization

if [ -z "$SUPERSTEP" ]
then
    python3 draw_graph.py --output_file=$OUTPUT_FILE
    deactivate
    ###############################
    echo "Opening graph image in browser..."
    # In MacOs:
    open graph.html
    # In Linux:
    # firefox graph.html
    # In Windows:
    # Explorer.exe graph.html
    ###############################
else
    python3 draw_graph.py --superstep=$SUPERSTEP
    deactivate
    ###############################
    echo "Opening graph image in browser..."
    # In MacOs:
    open graph-superstep-$SUPERSTEP.html
    # In Linux:
    # firefox graph.html
    # In Windows:
    # Explorer.exe graph.html
    ###############################
fi

cd ..
