cd ../src
go build -o pregel .
./pregel -type master -addr localhost -graph_file ../graphs/graph1.json