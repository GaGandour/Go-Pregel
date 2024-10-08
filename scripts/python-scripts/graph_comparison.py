import argparse
import json
import os

OUTPUT_GRAPH_FILE = "../../src/output_graphs/"
ANSWER_GRAPH_PATH = "../../answers/"


def generate_argparse() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Checks if the output graph is correct")
    parser.add_argument(
        "--verbose",
        help="Verbose mode",
        action="store_true",
    )
    parser.add_argument(
        "--algorithm",
        type=str,
        help="Output graph file",
        required=True,
    )
    return parser


def graph_name_from_file(file_name: str) -> str:
    return file_name.split("/")[-1]


def begin_test(graph_file: str) -> None:
    print("=========================================================")
    print(f"Testing output graph for {graph_name_from_file(graph_file)}...\n")


def print_answer_is_wrong(graph_file: str) -> None:
    print(f"[x] Output graph for {graph_name_from_file(graph_file)} is incorrect.")


def print_answer_is_correct(graph_file: str) -> None:
    print(f"[v] Output graph for {graph_name_from_file(graph_file)} is correct.")


def graph_tester(
    answer_file: str,
    pregel_file: str,
    verbose: bool = False,
) -> int:
    # This flag may change along with the function execution
    output_is_correct = True

    output_graph = {}
    answer_graph = {}
    with open(pregel_file, "r") as f:
        output_graph = json.load(f)
    with open(answer_file, "r") as f:
        answer_graph = json.load(f)

    answer_graph = answer_graph["Vertexes"]
    output_graph = output_graph["Vertexes"]

    answer_keys = answer_graph.keys()
    graph_keys = output_graph.keys()
    missing_keys = []
    for answer_key in answer_keys:
        if answer_key not in graph_keys:
            output_is_correct = False
            missing_keys.append(answer_key)
    extra_keys = []
    for graph_key in graph_keys:
        if graph_key not in answer_keys:
            output_is_correct = False
            extra_keys.append(graph_key)

    if not output_is_correct:
        if verbose:
            if extra_keys:
                print("\tThere are some extra vertices in your graph:")
                print(f"\tExtra vertices: {extra_keys}")
            if missing_keys:
                print("\tSome vertices are missing from your graph:")
                print(f"\tMissing vertices: {missing_keys}")
        return 0

    # check for correctness of the values:
    for key in answer_keys:
        output_vertex = output_graph[key]
        answer_vertex = answer_graph[key]

        output_value = output_vertex["Value"]
        answer_value = answer_vertex["Value"]
        for value_key in answer_value.keys():
            if value_key not in output_value:
                if verbose:
                    print(f'\tValue "{value_key}" is missing!')
                return 0
            if output_value[value_key] != answer_value[value_key]:
                if verbose:
                    print(f'\tValue "{value_key}" is incorrect!')
                return 0
    return 1

def test_graph(answer_file, pregel_file, verbose=False):
    begin_test(pregel_file)
    result = graph_tester(answer_file, pregel_file, verbose)
    if result == 0:
        print_answer_is_wrong(pregel_file)
    else:
        print_answer_is_correct(pregel_file)
    print("=========================================================")
    return result

if __name__ == "__main__":
    parser = generate_argparse()
    args = parser.parse_args()

    files_to_compare = os.listdir(ANSWER_GRAPH_PATH + args.algorithm + "/")

    correct = 0
    total = 0

    for file in files_to_compare:
        graph_file = args.algorithm + "/" + file
        answer_file = ANSWER_GRAPH_PATH + graph_file
        pregel_file = OUTPUT_GRAPH_FILE + graph_file
        total += 1
        correct += test_graph(answer_file, pregel_file, args.verbose)

    if correct == total:
        print(f"All {correct} of {total} graphs are correct!")
    else:
        print(f"{total - correct} of {total} graphs are incorrect!")
