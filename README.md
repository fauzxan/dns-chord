# dns-chord
Implementing DNS functionality using chord framework

🚀 [Problem description](https://github.com/fauzxan/dns-chord/blob/main/documentation/problem-description.md)

## Setup

### Local setup
To run locally, open number of terminals = number of nodes you want in the network. Then you need to run the following command:

```shell
go build && ./core -p <own port number | required> -u <some port number | optional>
```
#### Parameters:
##### -p
Refers to the port number that the client is going to run on, and listen to incoming connections from. This parameter is necessary in either cases - if you're the first node in the network, or not.

##### -u
This parameter is only necessary if you want to join an existing network. You need to know the port number of another node in the chord network, and use that to join the chord network. 

#### Environment variables:

```
IPADDRESS=127.0.0.1
```
Enter the IP address that you want your clients to run on. By default it is set to localhost. You may also expose your public IP address if you want to run accross different systems. 

### Docker setup
> Docker images are still in development phase!
To run docker container, just build docker image using 

```
    docker build --tag dns-chord-node .
```

If successfully built, then run

```
    docker run dns-chord-node
```
