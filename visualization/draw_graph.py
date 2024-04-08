import matplotlib.pyplot as plt
import networkx as nx
import json
import sys

VERTEX_VALUE_KEY = "Value"
EDGE_VALUE_KEY = "Value"
EDGES_KEY = "Edges"
EDGE_DESTINATION_KEY = "To"

NODE_SIZE = 800
ARROW_SIZE = 30
ARROW_STYLE = '-|>'


def vertex_value_to_display(vertex_id, value_dict):
    return f"""{vertex_id}::{value_dict["NextVertex"]}:{value_dict["Distance"]}"""

def edge_value_to_display(edge_id, value_dict):
    return value_dict["Cost"]

def print_graph_from_file(file_name):
    vertexes = {}
    # read from file
    with open(file_name, "r") as f:
        vertexes = json.load(f)
    
    # Create a directed graph
    G = nx.DiGraph()
    # Add nodes
    all_vertexes_ids = vertexes.keys()
    G.add_nodes_from(all_vertexes_ids)

    # Add directed edges with weights
    for vertex_id in vertexes:
        vertex = vertexes[vertex_id]
        edges = vertex.get(EDGES_KEY, [])
        for edge_id in edges:
            edge = edges[edge_id]
            edge_value = edge_value_to_display(edge_id, edge.get(EDGE_VALUE_KEY, None))
            if edge_value is not None:   
                G.add_edge(
                    vertex_id, 
                    edge.get(EDGE_DESTINATION_KEY), 
                    weight=edge_value,
                )
            else:
                G.add_edge(
                    vertex_id, 
                    edge.get(EDGE_DESTINATION_KEY)
                )

    # Generate positions for each node
    positions = nx.circular_layout(G, scale=6)  # Use a larger scale for more spread
    positions = nx.spring_layout(G, pos=positions, k=50, iterations=100)


    # Draw the directed graph
    nx.draw_networkx_nodes(G, positions, node_size=NODE_SIZE)
    nx.draw_networkx_edges(G, positions, arrowstyle=ARROW_STYLE, arrowsize=ARROW_SIZE)
    

    # Edge labels
    edge_labels = nx.get_edge_attributes(G, 'weight')

    # Draw edge labels
    nx.draw_networkx_edge_labels(G, positions, edge_labels=edge_labels)

    # Draw vertex labels
    node_labels = {vertex_id: vertex_value_to_display(vertex_id, vertexes[vertex_id].get(VERTEX_VALUE_KEY)) for vertex_id in all_vertexes_ids}
    node_labels = {k: v for k, v in node_labels.items() if v is not None}
    nx.draw_networkx_labels(G, positions, labels=node_labels, font_size=12)

    # Display the graph
    plt.axis('off')
    plt.show()

if __name__ == "__main__":
    file = sys.argv[1]
    if file:
        print_graph_from_file(file)

