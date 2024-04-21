# Go-Pregel



## How to use this

### Prerequisites

First, you will need Docker and Go installed. You can install Docker at [https://www.docker.com/](https://www.docker.com/) and Go at [https://go.dev/doc/install](https://go.dev/doc/install).

We use Go to execute the Pregel logic, and Docker to containerize it.

You'll also need Python. We'll use Python to visualize the output graph, after the Pregel algorithm has finished.

### Preparing the python environment

For our Python file to work, we'll need some libraries. I recommend to set up a virtual environment (venv) to contain the necessary libraries, although it's not mandatory. If you want to set up the venv, it's very simple:

```
python -m venv venv
```

Now, we must enter the virtual environment, install the libraries, and finally we can get out of the venv:
```
source venv/bin/activate     # enter venv
pip install -r requirements.txt # install
deactivate                   # exit venv
```

Your python environment is now ready. If you don't want to set the venv, just run `pip install requirements.txt`

### Generating Missing Files

To generate the files that you must fill up to use Pregel, run:
```
cd scripts
sh write_untracked_files.sh
```

### Writing the Algorithm

You can modify the following files:
+ `src/graph_package/user_defined_graph_methods.go`
+ `src/graph_package/user_defined_graph_types.go`
+ `src/graph_package/user_defined_utils_methods.go`
+ `src/graph_package/user_defined_utils_types.go`
+ `visualization/user_defined_value_displaying.py`

### Preparing the Docker image

This step has to be redone everytime you modify a Go file.

```
cd scripts
sh build_image.sh
```

### Running Pregel

Finally, we can run your algorithm in any graph in the `graphs/` folder.

```
sh start_docker.sh <NUMBER OF WORKERS> <Name OF GRAPH FILE>
```
In the end of the execution, a browser page with the output graph will open.