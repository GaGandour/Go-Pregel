import argparse
import json
import os

from pyvis.network import Network
from user_defined_value_displaying import edge_value_to_display, vertex_value_to_display

FILE = "../src/output_graphs/output_graph.json"
PREGEL_OUTPUT_DIR = "../src/output_graphs"
OUTPUT_FILE = "graph.html"

VERTEX_VALUE_KEY = "Value"
EDGE_VALUE_KEY = "Value"
EDGES_KEY = "Edges"
EDGE_DESTINATION_KEY = "To"


def generate_argparse() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Generate HTML output for the pregel graph")
    parser.add_argument(
        "--output_file",
        type=str,
        help="Output file name",
        default="",
    )
    parser.add_argument(
        "--superstep",
        type=int,
        help="Superstep number",
        default=-1,
    )
    return parser


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


def get_vertexes_from_args(output_file, superstep_number) -> dict:
    if output_file:
        full_output_file = os.path.join(PREGEL_OUTPUT_DIR, output_file)
        with open(full_output_file, "r") as f:
            return json.load(f)["Vertexes"]
    if superstep_number >= 0:
        vertexes = {}
        files_in_pregel_dir = [
            f for f in os.listdir(PREGEL_OUTPUT_DIR) if os.path.isfile(os.path.join(PREGEL_OUTPUT_DIR, f))
        ]
        filtered_files = [f for f in files_in_pregel_dir if f.startswith(f"SuperStep-{superstep_number}-")]

        for file in filtered_files:
            with open(os.path.join(PREGEL_OUTPUT_DIR, file), "r") as f:
                vertexes.update(json.load(f)["Vertexes"])
        global OUTPUT_FILE
        OUTPUT_FILE = "graph-superstep-" + str(superstep_number) + ".html"

        return vertexes
    return {}


if __name__ == "__main__":
    parser = generate_argparse()
    args = parser.parse_args()
    vertexes = get_vertexes_from_args(args.output_file, args.superstep)

    print_graph_from_dict(vertexes)
