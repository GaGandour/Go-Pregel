import argparse
import json
import subprocess
import sys

from pyvis.network import Network
from user_defined_value_displaying import edge_value_to_display, vertex_value_to_display


PREGEL_OUTPUT_DIR = "../src/output_graphs/"
PREGEL_OUTPUT_FILE = "../src/output_graphs/output_graph.json"
OUTPUT_FILE = "graph.html"

VERTEX_VALUE_KEY = "Value"
EDGE_VALUE_KEY = "Value"
EDGES_KEY = "Edges"
EDGE_DESTINATION_KEY = "To"


def generate_argparse() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Docker compose file generator for pregel algorithm")
    parser.add_argument(
        "--superstep",
        type=int,
        help="Superstep number to be visualized. If empty, the final output will be shown",
        default=-1,
    )
    parser.add_argument(
        "--pregel_output_file",
        type=str,
        help="Output file name",
        default=PREGEL_OUTPUT_FILE,
    )
    return parser


def get_superstep_logs(superstep: int):
    # Run the 'ls' command and capture the output
    result = subprocess.run(["ls", PREGEL_OUTPUT_DIR], capture_output=True, text=True, check=True)

    # Split the output into a list of strings
    files = result.stdout.splitlines()

    # Filter files with a certain prefix
    files = [file for file in files if file.startswith(f"SuperStep-{superstep}-")]

    return files


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
    parser = generate_argparse()
    args = parser.parse_args()
    graph = {}
    vertexes = {}
    temp_vertexes = {}
    if args.superstep == -1:
        # We are reading the final Pregel output
        with open(args.pregel_output_file, "r") as f:
            vertexes = json.load(f)["Vertexes"]
    else:
        files = get_superstep_logs(args.superstep)
        for file in files:
            with open(file, "r") as f:
                temp_vertexes = json.load(f)["Vertexes"]
            vertexes.update(temp_vertexes)
        OUTPUT_FILE = "graph-superstep-" + args.superstep + ".html"

    print_graph_from_dict(vertexes)
