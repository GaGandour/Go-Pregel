import json
import sys

from pyvis.network import Network
from user_defined_value_displaying import edge_value_to_display, vertex_value_to_display

FILE = "../src/output_graphs/output_graph.json"
OUTPUT_FILE = "graph.html"

VERTEX_VALUE_KEY = "Value"
EDGE_VALUE_KEY = "Value"
EDGES_KEY = "Edges"
EDGE_DESTINATION_KEY = "To"


def print_graph_from_dict(vertexes: dict):
    # Create a directed graph
    net = Network(notebook=True, cdn_resources="remote", select_menu=False, directed=True)

    # Add directed nodes
    for vertex_id in vertexes:
        vertex = vertexes[vertex_id]
        vertex_value = vertex_value_to_display(vertex_id, vertex.get(VERTEX_VALUE_KEY))
        net.add_node(
            vertex_id,
            label=vertex_value,
            physics=False,
        )

    # Add directed edges
    for vertex_id in vertexes:
        vertex = vertexes[vertex_id]

        edges = vertex.get(EDGES_KEY, [])
        for edge_id in edges:
            edge = edges[edge_id]
            edge_value = edge_value_to_display(edge_id, edge.get(EDGE_VALUE_KEY, None))
            if edge_value is not None:
                net.add_edge(
                    source=vertex_id,
                    to=edge.get(EDGE_DESTINATION_KEY),
                    title=edge_value,
                    arrowStrikethrough=False,
                )
            else:
                net.add_edge(
                    source=vertex_id,
                    to=edge.get(EDGE_DESTINATION_KEY),
                    arrowStrikethrough=False,
                )

    net.show(OUTPUT_FILE)


if __name__ == "__main__":
    vertexes = {}
    temp_vertexes = {}
    if len(sys.argv) == 2:
        print("This superstep doesn't exist")
        sys.exit(1)
    elif len(sys.argv) > 2:
        for arg in sys.argv[2:]:
            file = arg
            with open(file, "r") as f:
                temp_vertexes = json.load(f)
            vertexes.update(temp_vertexes)
        OUTPUT_FILE = "graph-superstep-" + sys.argv[1] + ".html"
    else:
        # read from file
        with open(FILE, "r") as f:
            vertexes = json.load(f)

    print_graph_from_dict(vertexes)
