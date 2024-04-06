import sys

def create_worker(worker_id: int) -> str:
    return f"""  pregel-worker-{worker_id}:
    image: pregel
    container_name: pregel-worker-{worker_id}
    command: ["./pregel", "-type", "worker", "-addr", "pregel-worker-{worker_id}", "-port", "5000{worker_id}", "-master", "pregel-master:5000"]
    ports:
      - "5000{worker_id}:5000{worker_id}"
    volumes:
      - FS:/lab_pregel
"""

def create_master(input_file) -> str:
    return f"""  pregel-master:
    image: pregel
    container_name: pregel-master
    command: ["./pregel", "-type", "master", "-addr", "pregel-master", "-graph_file", "{input_file}"]
    tty: true
    stdin_open: true
    ports:
      - "5000:5000"
    volumes:
      - FS:/lab_pregel
"""

def create_volumes() -> str:
    return f"""volumes:
  FS:
    external: true
"""

def create_docker_compose(num_workers: int, input_file: str) -> str:
    workers_description = "\n".join([create_worker(i) for i in range(1, num_workers+1)])
    return f"""version: '3'
services:
{create_master(input_file)}
{workers_description}
{create_volumes()}"""

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python write_docker_compose.py <num_workers> <graph_file>")
        sys.exit(1)
    print(create_docker_compose(int(sys.argv[1]), sys.argv[2]))