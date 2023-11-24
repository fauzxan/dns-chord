# dns-chord
Implementing DNS functionality using chord framework

🚀 [Problem description](https://github.com/fauzxan/dns-chord/blob/main/documentation/problem-description.md)

## Setup

### Local setup
To run locally, open number of terminals = number of nodes you want in the network. Then you need to run the following command:

```shell
go build && ./core
```

You will be asked to input the following:

1. Your current port number
2. Full IP address of the node you're using to join the network. If you are creating your own network, simply hit `ENTER` or `RETURN`

### Docker setup
> Docker images are still in development phase!
To run docker container, just build docker image using 

```shell
    docker build --tag dns-chord-node .
```

Build a docker volume called mydata (This is not needed anymore)
```shell
    docker volume create mydata
```

If successfully built, then run, as well as to bind the volume with the container, run 

```shell
    docker run -v mydata:/app/data  -it dns-chord-node
```
Do note that the -it tag is important to enable interactivity and also see colored output.
This mounts the "mydata" volume to the "/app/data" path inside the container.