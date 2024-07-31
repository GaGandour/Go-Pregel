#/bin/bash

SUPERSTEP=-1

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

cd ../..
source venv/bin/activate
cd visualization

if [ "$SUPERSTEP" -eq -1 ]; then
    python3 draw_graph.py
    deactivate
    open graph.html
else
    python3 draw_graph.py -superstep=$SUPERSTEP
    deactivate
    open graph-superstep-$SUPERSTEP.html
fi

cd ..
