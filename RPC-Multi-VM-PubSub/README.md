# RPC-Multi-VM-PubSub

This project implements a small distributed system using multiple virtual machines and network-based communication.

The main topics covered in this assignment are:

* Remote Procedure Call (RPC)
* Distributed services across multiple virtual machines
* Separation of service responsibilities
* User authentication through a remote service
* File and image serving
* Publish/Subscribe communication
* Memory usage monitoring and event notification

---

## Project Overview

The system consists of three independent virtual machines:

| Virtual Machine | Service               | Responsibility                                              |
| --------------- | --------------------- | ----------------------------------------------------------- |
| VM1             | Web Server            | Web interface, login, and communication with other services |
| VM2             | Authentication Server | User authentication and user data management                |
| VM3             | File Server           | Providing files and images                                  |

The Web Server communicates with the Authentication Server through RPC.

After successful authentication, the Web Server retrieves a file or image from the File Server over the network.

The Web Server also monitors its memory usage. When the memory consumption exceeds the configured threshold, an event is published and received by a subscriber.

---

## System Architecture

```text
                         User / Browser
                              |
                              v
                     +------------------+
                     |   VM1: Web       |
                     |     Server       |
                     +------------------+
                       /              \
                    RPC                HTTP/File
                     /                    \
                    v                      v
          +------------------+     +------------------+
          |   VM2: Auth      |     |   VM3: File      |
          |     Server       |     |     Server       |
          +------------------+     +------------------+

                     VM1: Web Server
                            |
                    Memory Monitoring
                            |
                            v
                    +------------------+
                    |    Publisher     |
                    |       |          |
                    |       v          |
                    |    Subscriber   |
                    +------------------+
                            |
                            v
                       Memory Alert
```

---

## Project Structure

```text
HW2/
├── README.md
├── report.pdf
│
├── rpc-study/
│   ├── rpc-summary.pdf
│   └── README.md
│
├── web-vm/
│   ├── main.go
│   ├── templates/
│   └── README.md
│
├── auth-vm/
│   ├── main.go
│   ├── users.json
│   └── README.md
│
├── file-vm/
│   ├── main.go
│   ├── files/
│   ├── images/
│   └── README.md
│
├── pubsub/
│   ├── publisher.go
│   ├── subscriber.go
│   └── README.md
│
└── screenshots/
```

---

## Components

### 1. RPC Study

The `rpc-study` directory contains the theoretical study of Remote Procedure Call.

The study covers:

* Definition and main concept of RPC
* Local vs. remote procedure calls
* Client and Server Stubs
* Serialization and Deserialization
* Interface Definition Language (IDL)
* Common RPC errors
* RPC vs. REST
* Comparison of different RPC technologies

The assignment requires comparison of at least three RPC technologies.

---

### 2. Web Server - VM1

The `web-vm` directory contains the Web Server running on VM1.

Its responsibilities include:

* Providing the login page
* Receiving username and password
* Sending authentication requests to VM2 through RPC
* Handling authentication results
* Displaying the appropriate response to the user
* Retrieving a file or image from VM3 after successful authentication
* Monitoring memory usage

The Web Server does not directly access the user database stored on VM2.

---

### 3. Authentication Server - VM2

The `auth-vm` directory contains the Authentication Server running on VM2.

Its responsibilities include:

* Storing user information
* Providing an RPC procedure for authentication
* Receiving username and password from VM1
* Returning the authentication result

The user data is stored on VM2 and is not directly accessed by the Web Server.

---

### 4. File Server - VM3

The `file-vm` directory contains the File Server running on VM3.

Its responsibilities include:

* Storing sample files and images
* Providing a network service for accessing those files
* Allowing VM1 to retrieve files or images over the network

The files are not stored locally on VM1.

---

### 5. Publish/Subscribe

The `pubsub` directory contains the Publish/Subscribe functionality.

The system monitors the memory usage of the Web Server.

When memory consumption exceeds the configured threshold:

```text
Web Server
    |
    | Memory usage exceeds threshold
    v
Publisher
    |
    | Event
    v
Subscriber
    |
    v
Memory Usage Alert
```

The published event contains information about the memory usage and configured threshold.

---

## Communication

The main communication paths in the system are:

```text
Browser
   |
   v
VM1 - Web Server
   |
   | RPC
   v
VM2 - Authentication Server
```

and:

```text
VM1 - Web Server
   |
   | Network Request
   v
VM3 - File Server
```

The Publish/Subscribe component is used separately for memory monitoring and event notification.

---

## Virtual Machines

The project requires at least three independent virtual machines.

### VM1 - Web Server

Runs the Web Server and provides the user interface.

### VM2 - Authentication Server

Runs the authentication service and stores user information.

### VM3 - File Server

Runs the file service and stores the files/images used by the Web Server.

Each service must communicate using the network addresses of the virtual machines rather than relying on `localhost`.

---

## Running the Project

Each component contains its own README with detailed instructions for installation, configuration, execution, and testing.

Please refer to:

```text
rpc-study/README.md
web-vm/README.md
auth-vm/README.md
file-vm/README.md
pubsub/README.md
```

Before running the complete system, make sure that the required virtual machines are running and can communicate with each other over the network.

---

## Recommended Startup Order

The general startup order is:

```text
1. Start VM2 - Authentication Server
2. Start VM3 - File Server
3. Start VM1 - Web Server
4. Start the Subscriber
5. Access the Web Server from the host/client
```

Detailed commands and configuration for each component are provided in the corresponding README files.

---

## Testing Scenario

The complete system can be tested using the following sequence:

### Authentication Test

1. Start the Authentication Server on VM2.
2. Start the Web Server on VM1.
3. Open the Web Server from the host or another client.
4. Enter invalid login credentials.
5. Verify that the login fails with an appropriate message.
6. Enter valid login credentials.
7. Verify that the login succeeds.

### File Server Test

After successful authentication:

1. Access the main page.
2. Request a file or image.
3. Verify that the Web Server retrieves the resource from VM3.
4. Verify that the resource is displayed successfully.

### Publish/Subscribe Test

1. Start the subscriber.
2. Start the Web Server.
3. Verify that memory usage is being monitored.
4. Increase the memory consumption of the Web Server using the provided mechanism.
5. Wait until the configured threshold is exceeded.
6. Verify that the Web Server publishes a memory event.
7. Verify that the subscriber receives the event.
8. Verify that a clear memory usage alert is displayed.

The assignment requires this scenario to be documented and supported by screenshots or terminal output.

---

## Error Handling

The system should handle common errors appropriately, including:

* Invalid login credentials
* Missing or invalid request parameters
* Network communication failures
* Authentication service unavailability
* File service unavailability
* Invalid RPC requests
* Memory monitoring and threshold-related errors

Specific error-handling details are documented in the README of each service.

---

## Technologies

The project uses technologies and concepts including:

* Go
* RPC
* HTTP
* JSON
* Virtual Machines
* Network Communication
* Publish/Subscribe
* Memory Monitoring

Possible RPC technologies specified by the assignment include:

* gRPC
* JSON-RPC
* XML-RPC
* HTTP/JSON-based RPC-like implementation

The selected RPC technology and the reason for choosing it are documented in the project report.

---

## Documentation

The project documentation is divided into several parts.

### Main Report

```text
report.pdf
```

The final report contains:

* Overall system architecture
* Role of each virtual machine
* IP addresses
* Selected RPC technology
* Reason for selecting the RPC technology
* Implemented procedures
* File transfer mechanism
* Publish/Subscribe mechanism
* Memory monitoring
* Testing results
* Screenshots and terminal outputs
* Implementation challenges and limitations

### Component Documentation

Each service contains a separate README describing how to run and test that component.

---

## Screenshots

The `screenshots/` directory contains screenshots or terminal outputs demonstrating the successful execution of the system.

Examples include:

* Virtual machine IP addresses
* Authentication service running
* File service running
* Web service running
* Successful login
* Failed login
* File/image retrieval
* Subscriber running
* Memory usage increase
* Memory alert

---

## Assignment Goals

This assignment demonstrates the following concepts:

* Understanding RPC and different RPC technologies
* Designing distributed services
* Running services across multiple virtual machines
* Remote authentication between independent services
* Separation of service responsibilities
* Network-based file access
* Publish/Subscribe event communication
* Monitoring system resources
* Documenting and testing a distributed system

---

## Notes

The project is designed to run in a Linux environment with multiple independent virtual machines.

All services should be configured so that their network communication is reproducible and clearly documented.

For detailed execution instructions, please refer to the README file inside each component directory.

