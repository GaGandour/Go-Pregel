def vertex_value_to_display(vertex_id, value_dict) -> str:
    """
    This function returns the text displayed under a vertex.

    params:
    - vertex_id: str. 
        A unique identifier for the vertex.
    - value_dict: dict.
        A dictionary with the exact same schema as in the VertexValue 
        struct inside /src/graph_package/user_defined_graph_types.go.
    """
    return f"""{vertex_id}"""


def edge_value_to_display(edge_id, value_dict) -> str:
    """
    This function returns the text displayed when hovering over an edge.

    params:
    - edge_id: str. 
        A unique identifier for the vertex.
    - value_dict: dict.
        A dictionary with the exact same schema as in the EdgeValue 
        struct inside /src/graph_package/user_defined_graph_types.go.
    """
    return f"""{edge_id}"""