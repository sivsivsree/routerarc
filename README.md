# RouterArc [Working on it] 
### API reverse proxy and Router Service.


RouterArc is a API gateway and reverse proxy service, can be used as an entry point of a microservice,
routing to diffrent services. 


## Goal (Initially) 🍺

### 1. Reverse Proxy  🔀

A reverse proxy is a server that sits in front of web servers and forwards client (e.g. web browser) requests to those web servers. Reverse proxies are typically implemented to help increase security, performance, and reliability. 
This is different from a forward proxy, where the proxy sits in front of the clients. With a reverse proxy, when clients send requests to the origin server of a website, those requests are intercepted at the network edge by the reverse proxy server. The reverse proxy server will then send requests to and receive responses from the origin server.

### 2. API Routing gateway  🚏

An API Gateway is a server that is the single entry point into the system. It is similar to the Facade pattern from object‑oriented design. The API Gateway encapsulates the internal system architecture and provides an API that is tailored to each client.

### 3. Loadbalancing 🚥

 A popular website that gets millions of users every day may not be able to handle all of its incoming site traffic with a single origin server. Instead, the site can be distributed among a pool of different servers, all handling requests for the same site. In this case, a reverse proxy can provide a load balancing solution which will distribute the incoming traffic evenly among the different servers to prevent any single server from becoming overloaded. In the event that a server fails completely, other servers can step up to handle the traffic.

<hr>

## Config File Format

The whole point of the project is to create a simple configuration based on json ```config.json```,
with all the features packed to run microservices 

```
router:[
  {
    servie: "/auth",
    upstream: [
      "http://localhost:8081",
      "http://localhost:8082",
    ],
    loadbalacer: "round-robin" //[round-robin, least-connection, iphash]
  },
  {
    servie: "/retrival",
    upstream: [
      "http://localhost:8084",
      "http://localhost:8085",
    ],
    loadbalacer: "round-robin" //[round-robin, least-connection, iphash]
  },  
],
proxy:[
    {
     from: "http://api.example.com",
     to: [
        "http://service1.example.com",
        "http://service2.example.com"
     ],
     loadbalacer: "round-robin" //[round-robin, least-connection, iphash]
   },
   {
     from: "http://example.com",
     to: [
        "http://localhost:8084",
        "http://service5.example.com"
     ],
     loadbalacer: "round-robin" //[round-robin, least-connection, iphash]
   }
]


```
