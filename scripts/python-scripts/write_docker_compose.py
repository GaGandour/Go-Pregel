import argparse

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
    parser.add_argument(
        "--checkpoint_frequency",
        type=int,
        help="Checkpoint frequency",
        default=-1,
    )
    return parser


def create_worker(worker_id: int, failure_step: int) -> str:
    cmd_elements = [
        "./pregel",
        "-type",
        "worker",
        "-addr",
        f"pregel-worker-{worker_id}",
        "-port",
        f"5000{worker_id}",
        "-master",
        f"pregel-master:{MASTER_PORT}",
    ]
    if failure_step >= 0 and worker_id == 1:
        cmd_elements.extend(
            [
                "-failure_step",
                f"{failure_step}",
            ],
        )

    cmd_string = ", ".join([f'"{element}"' for element in cmd_elements])
    cmd_string = f"[{cmd_string}]"

    return f"""  pregel-worker-{worker_id}:
    image: pregel
    container_name: pregel-worker-{worker_id}
    command: {cmd_string}
    ports:
      - "5000{worker_id}:5000{worker_id}"
    volumes:
      - ./graphs:/graphs
      - ./src/output_graphs:/src/output_graphs
"""


def create_master(
    input_file: str,
    debug: bool,
    num_workers: int,
    checkpoint_frequency: int,
) -> str:

    cmd_elements = [
        "./pregel",
        "-type",
        "master",
        "-port",
        f"{MASTER_PORT}",
        "-addr",
        "pregel-master",
        "-graph_file",
        f"/graphs/{input_file}",
    ]
    if debug:
        cmd_elements.append("-debug")
    if checkpoint_frequency > 0:
        cmd_elements.extend(
            [
                "-checkpoint_frequency",
                f"{checkpoint_frequency}",
            ],
        )
    cmd_string = ", ".join([f'"{element}"' for element in cmd_elements])
    cmd_string = f"[{cmd_string}]"
    
    master = f"""  pregel-master:
    image: pregel
    container_name: pregel-master
    command: {cmd_string}
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
    checkpoint_frequency: int,
) -> str:
    workers_description = "\n".join([create_worker(i, failure_step) for i in range(1, num_workers + 1)])
    return f"""services:
{create_master(input_file, debug, num_workers, checkpoint_frequency)}
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
            args.checkpoint_frequency,
        )
    )
