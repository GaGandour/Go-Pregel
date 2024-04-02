import sys

def create_worker(worker_id: int) -> str:
    return f"""  pregel-worker-{worker_id}:
    image: pregel
    container_name: pregel-worker-{worker_id}
    command: ["./pregel", "-type", "worker", "-addr", "pregel-worker-{worker_id}", "-port", "5000{worker_id}", "-master", "pregel-master:5000"]
    ports:
      - "5000{worker_id}:5000{worker_id}"
"""

def create_master() -> str:
    return """  pregel-master:
    image: pregel
    container_name: pregel-master
    command: ["./pregel", "-type", "master", "-addr", "pregel-master"]
    tty: true
    stdin_open: true
    ports:
      - "5000:5000"
"""

def create_docker_compose(num_workers: int) -> str:
    workers_description = "\n".join([create_worker(i) for i in range(1, num_workers+1)])
    return f"""version: '3'
services:
{create_master()}
{workers_description}"""

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python write_docker_compose.py <num_workers>")
        sys.exit(1)
    print(create_docker_compose(int(sys.argv[1])))