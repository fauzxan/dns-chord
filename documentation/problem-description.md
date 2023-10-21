This section will briefly discuss about the key features that are involved in chord and DNS systems, and also walk through the features that will be implemented in this project.

# Problem description

## Overall description

Traditional DNS architecture is inefficient and vulnerable and can lead to performance degradation, outages, and security vulnerabilities.

Traditional DNS uses a hierarchical structure, which serves requests from clients recursively or iteratively. The lookup times associated with such a system tend to be variable depending on potential caching. Also, the amount of load each level needs to take on is highly unequal. Moreover, the TLDs and upper-level servers are susceptible to DDoS attacks and prove to be a single point of failure.
Role of distributed systems in solving this issue
A distributed indexing system based on Chord would mitigate the issue of load balancing in DNS architectures if nodes in a distributed system can effectively look up domains without having to query vertically between different levels. The following properties of such a distributed indexing system built on Chord can be capitalized on:
- Resilience to single point of failure.
- Seamless node joins and departures.
- Scalable performance and faster lookup times.
- Elimination of hierarchy for decentralization improves load balancing.

In this project, we propose implementing **a chord-based domain query system that helps transition from a traditional DNS architecture to a distributed system** with the aforementioned benefits. 

## Key features planned

At its core, our project will aim to deliver the following features of the chord protocol:
- Lookup of domains that return the correct domain-IP mapping using any node in a distributed system.
- Account for failed nodes or new nodes in the system.
- Porting over of DNS records from DNS servers to a chord network. 
- <Add more>














<!--
# Ignore after this line



## Overview of chord
Chord is a dsitributed lookup protocol that stores a distributed hash table. Each node stores two types of key-value mappings:
1. **key-node mapping:** if the hash of the query is not directly present in the node, the node will find out the location of the node that contains that particular key.
2. **key-value mapping:** if the hash of the query is present in the node.

In steady state, each node stores at most O(log n) other nodes. However, there is a caveat in terms of performance when it comes to maintaining nodes that contain out-of-date information -> because it will be hard to maintain the consistency of the O(log n) state. 

### Properties of chord protocol:
Chord provides a distributed computiation of hash function. 
1. Chord assigns keys to nodes, by making use of **consistent hashing**.
2. **All nodes receive roughly the same number of keys**- this is a property of consistent hashing. This allows for effective load balancing. A given node will only receive the same number of queries as any other node over a long enough period of time. (We assume that no node contains keys that are more popular than the others)
3. **A given node doesn't contain routing information about all other nodes**. It only contains a small subset of this information; namely O(log n).

## Consistent hashing
We maintain an m-bit identifier for the following:
1. The nodes. This is obtained by hashing the nodes IP address. 
2. The keys that the nodes store. This is obtained by hashing the data.

Key k is assigned to the first node, whose identifier is the same as k's identifier, or its immediate successor node. (Denoted as successor(k)).
As an example, if a key has a hash value of 10, then it will either be assigned to the node, whose identifier is also 10, or the next node closest to the value of 10.

When a node joins the network, some of the keys that were previously assigned to the successor node n, will now be assigned to node n.
Similarly, when a node leaves the network, some of the keys that were previously assigned to it, will now be assigned to the successor of n. 

## Scalable lookup
One way to find a node that contains a given key is to hop from current node to successor node until you find the identifier of the node that is larger than the identifier

--!>
