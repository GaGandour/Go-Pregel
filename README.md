# Go-Pregel

## What is pregel?

Pregel is a distributed graph processing system developed by Google. It is designed to be scalable, fault-tolerant, and easy to program. The user defines a function that is executed by each vertex in the graph, and the system takes care of the rest.

This project is NOT Google's Pregel, but a simplified version of it, implemented in Go, based on the article "Pregel: A System for Large-Scale Graph Processing" by Grzegorz Malewicz, Matthew H. Austern, Aart J.C. Bik, James C. Dehnert, Ilan Horn, Naty Leiser, and Grzegorz Czajkowski.

The idea of this project is to implement a mini version of Pregel (called Go-Pregel), to understand how distributed graph processing models work and to learn how to use them.

This project offers 3 things:

### 1. A Graph Set.

Under the `graphs/` folder, you can find some graphs that you can use to test your Pregel algorithm. You can also create your own graph, following a similar format, and addapting the corresponding files (to be explained better in the next sections).

### 2. A Pregel Engine.

Our Go-Pregel is implemented under the `src/` folder. It uses a master-worker architecture, where the master is responsible for coordinating the workers, and the workers are responsible for executing the Pregel algorithm. The only thing you need to do to use Go-Pregel is to implement the Pregel logic in the `src/graph_package/user_defined_*.go` files. They are the only files you need to modify to use Go-Pregel. These contain the functions that will be used to determine the algorithm to be run. Depending on the algorithm you want to implement, the function will be different.

### 3. A Visualization Tool.

Inside the `visualization/` folder, you can find a Python script that will help you visualize the output of your Pregel algorithm. You can modify this script to display the information you want. The script will read the output of the Pregel algorithm and display it in a graph in your default browser. The display is interactive and you can move the nodes around.

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
source venv/bin/activate        # enter venv
pip install -r requirements.txt # install
deactivate                      # exit venv
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

In the next sections, we will explain what each file is for and how to use them.

### Running Pregel

Finally, we can run your algorithm in any graph in the `graphs/` folder.

```
sh start_docker.sh <NUMBER OF WORKERS> <Name OF GRAPH FILE>
```
In the end of the execution, a browser page with the output graph will open.

## Writing the Algorithm

### The Graph format

All of the following types can be defined by the user in the `src/graph_package/user_defined_graph_types.go` file.

The graphs are represented here as a dictionary. Each key is a vertex's ID (of type `VertexIdType`, defined by the user), and the value represents the vertex itself. Each vertex must have a `Value` field and a `Edges` field. The `Value` field is the value of the vertex, and its type depends on the algorithm and is defined by the user (in the class `VertexValue`), and the `Edges` field is another dictionary with the edges that come out of the vertex. You can give whatever key name you want to your edges, but they need to be unique. Each edge has a`To` field, of the type `VertexIdType`, and a `Value` field, of the type `EdgeValue`, defined by the user. The `To` field is the ID of the vertex that the edge goes to, and the `Value` field is the value of the edge. This `EdgeValue` can be useful or not, it depends on the algorithm. This can be used, for example, to store the weight of the edge. If you don't want to use the `EdgeValue` in your algorithm, just leave it as an empty struct. The same thing goes for the `VertexValue`.

### The Graph methods

#### Understanding Go-Pregel: A brief explanation

To understand how to write the algorithm, you need to understand how the Pregel framework (or, in this case, the Go-Pregel framework) works. The Pregel framework is based on the concept of supersteps. In each superstep, each vertex can do the following actions:

1. Read messages sent by other vertices in the previous superstep.
2. Change its own value, depending on its own value and on the messages it received and on the values of its outgoing edges. Please note that a vertex does NOT have access to the values of its neighbors, unless they send messages with this information.
3. Send messages to other vertices (not necessarily neighbors, the only requirement is that the ID of the target vertex is known). This message will only be read by the target vertex only in the next superstep.
4. Vote to halt. If all vertices vote to halt, the algorithm stops. If a vertex votes to halt, it stops executing the actions 2 and 3 in the next supersteps unless it is activated. A vertex is activated again if and only if it receives a message from another vertex.

#### The Pregel Message

Now that you know what a message is, you can write the `PregelMessage` struct in the file `src/graph_package/user_defined_graph_types.go`. This struct is used to send messages between vertices. Feel free to organize the struct as you wish. It should reflect the implementation of your algorithm.

#### The Methods

The user must implement a total of 2 to 3 methods (the last one is optional) in the `src/graph_package/user_defined_graph_methods.go` file. These methods are:

1. `vertex.ComputeInSuperStepZero`: This is the compute function used in the first superstep for each vertex. If it's not necessary, you can simply call `vertex.Compute([])` in this function.
2. `vertex.Compute`: This is the compute function used in the other supersteps for each vertex. This function receives a slice containing the messages sent to the vertex in the previous superstep. The user must implement the logic of the algorithm here.
3. `CombinePregelMessages`: This function is totally optional and its only effect is a slight improvement on the network traffic when executing pregel in large scale graphs, and only in certain algorithms. If you don't want to implement it, you can just return the only argument it receives: a list of messages. A deeper explanation of this function and what it does is given in the docstring of the function.

In the `Compute` method, we usually follow this logic, but adapting it according to the algorithm:

1. Interpret/Read the messages received. The parameter `messageList` contains the messages sent to the vertex in the previous superstep. You should read those messages and decide what to do with them. It's like starting a day of work: you read your emails and decide what to do in your day.
2. Decide what to do with your onw value (the vertex's value). You can change it or not, depending on the algorithm and on the messages received. You may want to do 2. at the same time as 1. If the algorithm envolves finding the maximum value of all vertices, for example, you could pass your current value as a message to other vertices and change your own value if you receive a message with a higher value than yours. In this case, you would be doing 1. and 2. at the same time.
3. Send messages to other vertices. You can send messages to any vertex you want, not necessarily neighbors. You can also send messages to yourself. You can send as many messages as you want. The only requirement is that the ID of the target vertex is known. The content of your message is defined by the `PregelMessage` struct. The function you want to use to send messages is `vertex.PrepareMessageToVertex`, which receives the ID of the target vertex and the content of the message. The reason this method is called `PrepareMessageToVertex` is because the message is not sent immediately. It is stored in a buffer and sent only at the end of the superstep.
4. Vote to halt. If you don't want to do anything in the next supersteps, you can vote to halt. If all vertices vote to halt, the algorithm stops. If you want to do something in the next supersteps, you can simply not vote to halt. The method to do this is `vertex.VoteToHalt`. The method doesn't receive any arguments.

### Graph Visualization Methods

Your Go-Pregel algorithm is ready to go. But when running it, the UI can be a bit confusing. That's why you can customize the graph info visualization in the `visualization/user_defined_value_displaying.py` file. This file is responsible for displaying the information of the vertices and edges in the graph. You can change the way the information is displayed. The file is quite simple and there are docstrings indicating the purpose of each function and what each argument means.

## Testing

<!-- TODO: There must be a way to test some algoritms, like big unit tests -->
