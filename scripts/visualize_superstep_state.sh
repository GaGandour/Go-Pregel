#/bin/bash

for arg in "$@"
do
    case $arg in
        -h|--help)
        echo "Usage: ./visualize_superstep_state.sh -superstep=<superstep number>"
        exit 0
        ;;
    esac
    case $arg in
        -superstep=*)
        SUPERSTEP="${arg#*=}"
        shift
        ;;
    esac
done

if [ -z "$SUPERSTEP" ]
  then
    echo "Missing SUPERSTEP. Run ./visualize_superstep_state.sh with -h or --help for more information on the necessary arguments."
    exit 1
fi

cd ..
source venv/bin/activate
cd visualization
python3 draw_graph.py $(ls ../src/output_graphs/SuperStep-$SUPERSTEP*)
deactivate
open graph.html
cd ..