import json
from pyvis.network import Network

FILE = "../src/output_graphs/output_graph.json"

VERTEX_VALUE_KEY = "Value"
EDGE_VALUE_KEY = "Value"
EDGES_KEY = "Edges"
EDGE_DESTINATION_KEY = "To"


def vertex_value_to_display(vertex_id, value_dict):
    return f"""{vertex_id}\n{value_dict["NextVertex"]}:{value_dict["Distance"]}"""

def edge_value_to_display(edge_id, value_dict):
    return value_dict["Cost"]

def print_graph_from_file(file_name):
    vertexes = {}
    # read from file
    with open(file_name, "r") as f:
        vertexes = json.load(f)
    
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

    net.show("graph.html")

if __name__ == "__main__":
    print_graph_from_file(FILE)

