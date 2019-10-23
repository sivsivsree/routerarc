# RouterArc [Working on it] 
### API reverse proxy and Router Service.
[![Actions Status](https://github.com/sivsivsree/routerarc/workflows/Build/badge.svg)](https://github.com/sivsivsree/routerarc/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/sivsivsree/routerarc)](https://goreportcard.com/report/github.com/sivsivsree/routerarc)

RouterArc is an API gateway and reverse proxy service, can be used as an entry point of a microservice,
routing to diffrent services. 


## Goal (Initially) 🍺

### 1. Reverse Proxy  🔀

A reverse proxy is a server that sits in front of web servers and forwards client (e.g. web browser) requests to those web servers. Reverse proxies are typically implemented to help increase security, performance, and reliability. 
This is different from a forward proxy, where the proxy sits in front of the clients. With a reverse proxy, when clients send requests to the origin server of a website, those requests are intercepted at the network edge by the reverse proxy server. The reverse proxy server will then send requests to and receive responses from the origin server.

### 2. API Routing gateway  🚏

An API Gateway is a server that is the single entry point into the system. It is similar to the Facade pattern from object‑oriented design. The API Gateway encapsulates the internal system architecture and provides an API that is tailored to each client.

### 3. Loadbalancing 🚥

 A popular website that gets millions of users every day may not be able to handle all of its incoming site traffic with a single origin server. Instead, the site can be distributed among a pool of different servers, all handling requests for the same site. In this case, a reverse proxy can provide a load balancing solution which will distribute the incoming traffic evenly among the different servers to prevent any single server from becoming overloaded. In the event that a server fails completely, other servers can step up to handle the traffic.


### 4. Static File Serving
  - Can be directly attached to a port
  > Note: if you have and angular application and want to port forward the build HTML and assets to a server u can use the proxy serve 
<hr>


## Usage

To run the service use, 

```
   routerarc -config=rules.json  
```

 
 
##### Flags :

- `` -h `` : to view all the command line arguments. 
- `` -config=<filename>`` : to specify the configurations.


<br>

<br>

  #### Using Docker 
 
 > The routing configuration file 'rules' should be in the volume. 
 > ``/var/routerarc`` in the below docker config.
 
 <br>
 
 ``` docker run --name somevol -v /var/routerarc:/rules routerarch ```
 
 
 
 <br>

## Config File Format

The whole point of the project is to create a simple configuration based on json ```config.json```,
with all the features packed to run microservices 

> The configuration can either be in YAML format or in JSON format.
> If no configuration file is provided it will look for ```rules.yaml```

 ##### 1. YAML configuration Example ```config.yaml```

```
router:
  - port: '8080'
    case:
      - service: "/auth"
        loadbalacer: round-robin
        upstream:
          - http://localhost:8081
          - http://localhost:8082
          -
      - servie: "/mobile/auth"
        upstream:
          - http://localhost:8086
          - http://localhost:8045
        loadbalacer: round-robin

      - servie: "/retrival"
        upstream:
          - http://localhost:8084
          - http://localhost:8085
        loadbalacer: round-robin
proxy:
  - name: backend
    port: '8000'
    to:
      - https://api.github.com
      - https://google.com
    loadbalacer: round-robin

  - name: frontend
    port: '9000'
    static: "public"
    
```

<br>

 ##### 2. JSON configuration Example ```rules.yaml```
```
{
  "router": [
    {
      "port": "8080",
      "case": [
        {
          "service": "/auth",
          "loadbalacer": "round-robin",
          "upstream": [
            "http://localhost:8081",
            "http://localhost:8082"
          ]
        },
        {
          "servie": "/mobile/auth",
          "upstream": [
            "http://localhost:8086",
            "http://localhost:8045"
          ],
          "loadbalacer": "round-robin"
        },
        {
          "servie": "/retrival",
          "upstream": [
            "http://localhost:8084",
            "http://localhost:8085"
          ],
          "loadbalacer": "round-robin"
        }
      ]
    }
  ],
  "proxy": [
    {
      "name": "backend",
      "port": "8000",
      "to": [
        "http://service1.ae",
        "http://service2.example.com"
      ],
      "loadbalacer": "round-robin"
    },
    {
      "name": "frontend",
      "port": "9000",
      "to": [
        "https://jsonplaceholder.typicode.com",
        "http://example.com"
      ],
      "loadbalacer": "round-robin"
    }
  ]
}


```
