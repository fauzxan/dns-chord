This section will briefly discuss about the key features that are involved in chord and DNS systems, and also walk through the features that will be implemented in this project.

# 1. Problem description
One of the essential features of the modern Internet is the domain name system (DNS), which is a
naming system that associates human-readable domain names to their corresponding numerical IP
addresses. However, traditional DNS architecture may have various performance inefficiencies and
vulnerabilities regarding security and outages. DNS implements a hierarchical structure, which can
serve requests from clients either recursively or iteratively. The lookup times associated with such a
system tend to be variable depending on potential caching, and the amount of load each level needs to
take on is highly unequal, leading to problems associated with centralization. For instance, the TLDs
and upper-level servers can become a single point of failure and make it susceptible to DDoS attacks.

# 2. Role of Distributed Systems in tackling this issue

A distributed indexing system based on Chord would mitigate the issue of load balancing in DNS
architectures if nodes in a distributed system can effectively look up domains without having to query
vertically between different levels.

The following properties of such a distributed indexing system built on Chord can be capitalized on:
- Resilience to a single point of failure.
- Seamless node joins and departures.
- Scalable performance and faster lookup times.
- Elimination of hierarchy for decentralization improves load balancing.

In this project, we propose implementing a chord-based domain query system that helps transition
from a traditional DNS architecture to a distributed system with the aforementioned benefits.

# 3. Key features

At its core, our project will aim to deliver the following features of the chord protocol:
- **Correctness**: Lookup of domains that return the correct domain-IP mapping using any node in a distributed system.
- **Availability and Fault Tolerance**: Account for failed nodes or new nodes in the system.
- **Porting over of DNS records** from DNS servers to a chord network.
- **Scalability**: It should efficiently provide mapping efficiently given a large number of domain-IP pairs in the hash table.
- **Decentralized**: Fully distributed with no one node more important than the other

# 4. Implementation plan
In this project, we will be implementing the Chord protocol as described in this paper. We will be
using the ‘Go’ Programming Language in our implementation of the protocol.
The protocol described in the paper will be completely implemented from scratch using the help of a
few packages that are part of the ‘Go’ Language Package Suite like the following:

- **Net Package**: The Net Package will be useful to interface with TCP/IP Sockets and DNS lookups in the Application.
- **Crypto Package**: The Crypto Package allows us to use the built-in SHA1 package for achieving consistent hashing as part of the Chord Protocol.
- **Cobra Package**: The Cobra Package will allow us to build a clean and modern command line interface that will help us log the changes in the network.

The functionality from these packages will be used to achieve a functional/working version of the
chord protocol.

# 5. Validation plan

As described above, our tool is a means to transition from the issues-ridden legacy DNS to the faster,
decentralized, and more fault-tolerant Chord-based DNS architecture. Validation for this system can
translate into 2 use cases:

## 5.1 Application use cases
- **DNS Lookup**: Demonstrate the ability of our system to effectively lookup domain-IP mappings.
If the domain name resides on the Chord network, the lookup call will fetch the DNS record from
the network. If it doesn’t exist on the network yet, a lookup will be performed on the legacy
DNS servers. Subsequently, the record will be inserted with its own key into the Chord network.
- **Node joins and departures**: Simulate scenarios where new nodes (Chord network participants)
join, or nodes leave (either voluntarily or due to failure). Showcase how our system can dynamically
adapt to such changes and maintain its functionality and correctness.


## 5.2 Testing
- **Unit testing**: To develop unit tests for different components of our Chord network. These tests
will account for scenarios such as lookups, nodes joining as well and handling failures.
- **Integration testing**: To develop system tests to ensure various components work together without
failing. Deal with edge cases that arise due to the integration of different components. For
example, the DNS record lookup function should work seamlessly even after the departure of
nodes from the network.

## 5.3 Simulation of Distributed Nodes
We will use separate computers to act as nodes to participate in the network.
- Simulate 3-5 nodes on a network residing on the Local Area Network (LAN). Perform functions
to demonstrate the use cases mentioned above. Use Terminal logs to show the correctness of the
functions.
- Make DNS requests for new records. Demonstrate how a call to the legacy DNS servers is followed
by inserting the record into the Chord network. Demonstrate that future lookups for the same
DNS record are routed to the Chord network.
- Simulate a node leaving due to voluntary departure. Demonstrate that the DNS record previously
hosted by this node still exists in the Chord network.
- Simulate a node joining the network. Demonstrate load balancing of DNS records; some records
are shifted from the original node to the new node.

If time permits, we may explore:
- The use of scripts to spin up docker containers and simulate these virtual nodes performing the
various functions.
- A visual depiction of nodes joining and leaving and the DNS records that reside on each participating
node.
