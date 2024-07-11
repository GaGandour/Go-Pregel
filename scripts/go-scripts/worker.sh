cd ../src
go build -o pregel .
./pregel -type worker -addr localhost -port 5000$1 -master localhost:5000