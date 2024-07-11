import argparse
import sys

MASTER_PORT = 50000


def generate_argparse() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Docker compose file generator for pregel algorithm")
    parser.add_argument(
        "--num_workers",
        type=int,
        help="Number of workers",
        required=True,
    )
    parser.add_argument(
        "--graph_file",
        type=str,
        help="Input graph file for pregel algorithm",
        required=True,
    )
    parser.add_argument(
        "--debug",
        action="store_true",
        help="Debug mode",
    )
    parser.add_argument(
        "--failure_step",
        type=int,
        help="Failure step",
        default=-1,
    )
    # parser.add_argument(
    #     "--checkpoint_frequency",
    #     type=int,
    #     help="Checkpoint frequency",
    #     default=-1,
    # )
    return parser


def create_worker(worker_id: int, failure_step: int) -> str:
    command = f"""["./pregel", "-type", "worker", "-addr", "pregel-worker-{worker_id}", "-port", "5000{worker_id}", "-master", "pregel-master:{MASTER_PORT}\""""
    if failure_step >= 0 and worker_id == 1:
        command += f""", "-failure_step", "{failure_step}\""""
    command += "]"

    return f"""  pregel-worker-{worker_id}:
    image: pregel
    container_name: pregel-worker-{worker_id}
    command: {command}
    ports:
      - "5000{worker_id}:5000{worker_id}"
    volumes:
      - ./graphs:/graphs
      - ./src/output_graphs:/src/output_graphs
"""


def create_master(input_file: str, debug: bool, num_workers: int) -> str:
    master = f"""  pregel-master:
    image: pregel
    container_name: pregel-master
    command: ["./pregel", "-type", "master", "-port", "{MASTER_PORT}", "-addr", "pregel-master", "-graph_file", "/graphs/{input_file}"{', "-debug"' if debug else ""}]
    tty: true
    stdin_open: true
    ports:
      - "{MASTER_PORT}:{MASTER_PORT}"
    volumes:
      - ./graphs:/graphs
      - ./src/output_graphs:/src/output_graphs
    depends_on:
"""
    for i in range(1, num_workers + 1):
        master += f"      - pregel-worker-{i}\n"
    return master


def create_volumes() -> str:
    return """volumes:
  FS:
    external: true
"""


def create_docker_compose(
    num_workers: int,
    input_file: str,
    debug: bool,
    failure_step: int,
) -> str:
    workers_description = "\n".join([create_worker(i, failure_step) for i in range(1, num_workers + 1)])
    return f"""version: '3'
services:
{create_master(input_file, debug, num_workers)}
{workers_description}
"""


if __name__ == "__main__":
    parser = generate_argparse()
    args = parser.parse_args()
    print(
        create_docker_compose(
            args.num_workers,
            args.graph_file,
            args.debug,
            args.failure_step,
        )
    )
