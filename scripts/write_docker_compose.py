import sys

MASTER_PORT = 50000


def create_worker(worker_id: int) -> str:
    return f"""  pregel-worker-{worker_id}:
    image: pregel
    container_name: pregel-worker-{worker_id}
    command: ["./pregel", "-type", "worker", "-addr", "pregel-worker-{worker_id}", "-port", "5000{worker_id}", "-master", "pregel-master:{MASTER_PORT}"]
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
    return f"""volumes:
  FS:
    external: true
"""


def create_docker_compose(num_workers: int, input_file: str, debug: bool) -> str:
    workers_description = "\n".join([create_worker(i) for i in range(1, num_workers + 1)])
    return f"""version: '3'
services:
{create_master(input_file, debug, num_workers)}
{workers_description}
"""


if __name__ == "__main__":
    num_args = len(sys.argv) - 1
    if num_args != 2 and num_args != 3:
        print("Usage: python write_docker_compose.py <num_workers> <graph_file> [debug]")
        sys.exit(1)
    print(create_docker_compose(int(sys.argv[1]), sys.argv[2], num_args == 3))
