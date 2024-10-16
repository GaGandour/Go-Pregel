# Go-Pregel

## What is pregel?

Pregel is a distributed graph processing system developed by Google. It is designed to be scalable, fault-tolerant, and easy to program. The user defines a function that is executed by each vertex in the graph, and the system takes care of the rest.

This project is NOT Google's Pregel, but a simplified version of it, implemented in Go, based on the article "Pregel: A System for Large-Scale Graph Processing" by Grzegorz Malewicz, Matthew H. Austern, Aart J.C. Bik, James C. Dehnert, Ilan Horn, Naty Leiser, and Grzegorz Czajkowski.

### A Super-Oversimplified Explanation of Pregel

The objective of Pregel is to process a graph in a distributed way. The graph is divided into partitions, and each partition is processed by a worker. However, each vertex acts as if it is a single node, even if it shares a worker with many other vertices (which is usually the case).

The Pregel framework is based on the concept of supersteps. In each superstep, each vertex can do the following actions:

1. Read messages sent to it by other vertices in the previous superstep.
2. Change its own value, depending on its own value and on the messages it received and on the values of its outgoing edges. Please note that a vertex does NOT have access to the values of its neighbors, unless they send messages with this information.
3. Send messages to other vertices (not necessarily neighbors, the only requirement is that the ID of the target vertex is known). This message will only be read by the target vertex, and only in the next superstep.
4. Vote to halt. If all vertices vote to halt, the algorithm stops. If a vertex votes to halt, it stops executing the actions 2 and 3 in the next supersteps unless it is activated. A vertex is activated again if and only if it receives a message from another vertex.

## What is this project?

The idea of this project is to implement a mini version of Pregel in Golang (called Go-Pregel), to understand how distributed graph processing models work and to learn how to use them.

This project offers 3 things:

### 1. A Graph Set.

Under the `graphs/` folder, you can find some graphs that you can use to test your Pregel algorithm. You can also create your own graph, following a similar format, and addapting the corresponding files (to be explained better in the next sections).

### 2. A Pregel Engine.

Our Go-Pregel is implemented under the `src/` folder. It uses a master-worker architecture, where the master is responsible for coordinating the workers, and the workers are responsible for executing the Pregel algorithm. The only thing you need to do to use Go-Pregel is to implement the Pregel logic in the `src/graph_package/user_defined_*.go` files. They are the only files you need to modify to use Go-Pregel. These contain the functions that will be used to determine the algorithm to be run. Depending on the algorithm you want to implement, the function will be different.

### 3. A Visualization Tool.

Inside the `visualization/` folder, you can find a Python script that will help you visualize the output of your Pregel algorithm. You can modify this script to display the information you want. The script will read the output of the Pregel algorithm and display it in a graph in your default browser. The display is interactive and you can move the nodes around.

## How to install and use Go-Pregel

### Prerequisites

First, you will need Docker and Go installed. You can install Docker at [https://www.docker.com/](https://www.docker.com/) and Go at [https://go.dev/doc/install](https://go.dev/doc/install).

We use Go to execute the Pregel logic, and Docker to containerize it.

You'll also need Python. We'll use Python (you can download it at [https://www.python.org/downloads/](https://www.python.org/downloads/)) to visualize the output graph, after the Pregel algorithm has finished. Python is also used to write the `docker-compose.yml` file, which is used to run the Pregel algorithm.

### If you are using Windows

If you are not using Windows, you can skip this section.

If you are using windows, you'll need wsl to run pregel. You can install it on Windows 10 by following the instructions at [https://docs.microsoft.com/en-us/windows/wsl/install-win10](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

After installing wsl, you'll need to run wsl and run the following commands:

```bash
sudo apt update
sudo apt install python3-pip
sudo apt install python3-venv
```

This guarantees that you'll be able to create a virtual environment on wsl.

Furthermore, you'll need to replace every appearence of the word `open` in the files `./scripts/execution/start_pregel.sh` and `./scripts/execution/visualize_superstep_state.sh` to `Explorer.exe`. This is because the `open` command is not available on wsl.

### If you are using Linux

If you are not using Linux, you can skip this section.

If you are using Linux, you'll need to replace every appearence of the word `open` in the files `./scripts/execution/start_pregel.sh` and `./scripts/execution/visualize_superstep_state.sh` to `firefox` or `google-chrome`, depending on your preferences. This is because the `open` command is not available on Linux.

### Preparing the python environment

For our Python file to work, we'll need some libraries. It is recommended to set up a virtual environment (venv) to contain the necessary libraries, although it's not mandatory. If you want to set up the venv, it's very simple:

```
python -m venv venv
```

Depending on the machine, you'll might need to use `python3` instead of `python`.

Now, we must enter the virtual environment, install the libraries, and finally we can get out of the venv:

```bash
source venv/bin/activate        # enter venv
pip install -r requirements.txt # install
deactivate                      # exit venv
```

Your python environment is now ready. If you don't want to set the venv, just run `pip install -r requirements.txt`. Be aware that, with this, the libraries will be installed in your global python environment.

### Generating Missing Files

To generate the files that you must fill up to use Pregel, run:
```bash
# Run from the root of the project
cd scripts/prepare-repo
./write_untracked_files.sh
cd ../..
```

### Writing the Graph Algorithm

You can modify the following files:
+ `src/graph_package/user_defined_graph_methods.go`
+ `src/graph_package/user_defined_graph_types.go`
+ `src/graph_package/user_defined_utils_methods.go`
+ `src/graph_package/user_defined_utils_types.go`
+ `visualization/user_defined_value_displaying.py`

To understand how to write a pregel algorithm, read the [pregel_writing_guide.md](https://github.com/GaGandour/Go-Pregel/blob/main/pregel_writing_guide.md) file.

### Running Pregel

Finally, we can run your algorithm in any graph in the `graphs/` folder by using the `start_docker.sh` script under the `/scripts/execution/` folder. Before using it, you MUST be in this `/scripts/execution/` folder before running it. To see the usage of the script and/or see the available arguments to pass to it, run `./start_docker.sh -h` or `./start_docker.sh --help`. **Remember to start the docker deamon!**

```bash
# Run from the root of the project
cd scripts/execution
./start_pregel.sh -h
cd ../..
```

In the most simple use case, the command will look like this, where the graph_file is the path to the graph file relative from the `graphs/` folder:

```bash
# Run from the root of the project
cd scripts/execution
./start_pregel.sh -num_workers=<number_of_workers> -graph_file=<graph_file>
cd ../..
```

In the end of the execution, a browser page with the output graph will open.

### Visualizing the Output

If you want to visualize the output of the Pregel algorithm without running it again, or if you want to visualize the Pregel state in a certain superstep (or even the initial state), you can use the `visualize_superstep_state.sh` script inside `./scripts/execution/`. The usage is similar to the `start_pregel.sh` script. If you want to visualize a superstep, you run:

```bash
# Run from the root of the project
cd scripts/execution
./visualize_superstep_state.sh -superstep=<superstep>
cd ../..
```

However, if you want to visualize the final pregel outputfor a certain graph, you can run the following commands, where the output file is also the relative path to the graph file from the `graphs/` folder:

```bash
# Run from the root of the project
cd scripts/execution
./visualize_superstep_state.sh -output_file=<output_file>
cd ../..
```

### Running tests for a certain algorithm

After writing a Pregel algorithm, you can test it in every available graph inside `./graphs/<algorithm>/` folder, by using the `./scripts/execution/test_pregel_algorithm.sh` script. To understand the flags and arguments that you can pass to the script, run:

```bash
# Run from the root of the project
cd scripts/execution
./test_pregel_algorithm.sh -h
cd ../..
```

### Debugging

When running Go-Pregel in a graph, besides the `-debug` option which logs the graph state in every superstep, when running `./scripts/execution/start_pregel.sh`, the logs from each worker are written in the `./docker_logs/` folder. It can be useful to debug your algorithm.
