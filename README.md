# dns-chord
Implementing DNS functionality using chord framework

## Setup
Just run go run main.go to run without docker container. 

To run docker container, just build docker image using 

```
    docker build --tag dns-chord-node .
```

If successfully built, then run

```
    docker run dns-chord-node
```