# Faraway Client

This is a Go client application that connects to a TCP server, handles proof-of-work challenges if required, and processes the response.

## Prerequisites

- Docker installed on your machine.
- The TCP server is running and accessible.

## Building the Docker Image

To build the Docker image for the Go client application, navigate to the project directory and run the following command:

```sh
docker build -t faraway-client -f ./scripts/docker/client.dockerfile .
docker run --network faraway-network -e SERVER_ADDRESS=faraway-server:8080 faraway-client
```

# Faraway Server

This is a TCP server application built with Go. It handles client connections, implements a proof-of-work mechanism for DDoS protection, and returns quotes.

## Prerequisites

- Docker installed on your machine.

## Building the Docker Image

To build the Docker image for the server application, navigate to the project directory and run the following command:

```sh
docker build -t faraway-server .
docker run -p 8080:8080 -e DDOS_PROTECTION=HASHCASH -e PORT=8080 -e READ_TIMEOUT=750ms -e WRITE_TIMEOUT=750ms faraway-server
```

# Faraway Server Protocol

This document describes the protocol used by the Faraway Server, which listens for TCP connections on port 8080 and implements a proof-of-work mechanism for DDoS protection.

## Protocol Overview

1. **Client Connection**: A client establishes a TCP connection to the server on port 8080.
2. **Resource Request**: The client sends the resource name "word_of_wisdom".
3. **Proof-of-Work Challenge**: The server responds with a proof-of-work challenge in the form of a Hashcash header.
4. **Challenge Solution**: The client solves the challenge and sends the solution back to the server.
5. **Quote Response**: Upon verification of the solution, the server sends a quote back to the client.

## Communication Details

### Client Request

- **Resource Request**: The client sends the resource name "word_of_wisdom" followed by a newline character (`\n`).



### Server Response

- **Proof-of-Work Challenge**: The server responds with a Hashcash header in the following format:

X-Hashcash: <version>:<date>:<resource>::<complexity>:<nonce>\n

Example: X-Hashcash: 1:20240618:word_of_wisdom::5:7c72bb9946334bd1a776a405671a74aa


- `version`: Hashcash version, always `1`.
- `date`: The date in `YYYYMMDD` format.
- `resource`: The requested resource, in this case, `word_of_wisdom`.
- `complexity`: The number of leading zero bits required in the solution.
- `nonce`: A random string used in the challenge.

### Client Solution

- **Challenge Solution**: The client calculates a solution to the Hashcash challenge and sends it back to the server, followed by a newline character (`\n`).


### Server Verification and Quote Response

- **Quote Response**: If the solution is verified successfully, the server sends a quote back to the client.


## Example Communication Flow

1. **Client sends resource request**:
  ```
  word_of_wisdom\n
  ```

2. **Server sends proof-of-work challenge**:
  ```
  X-Hashcash: 1:20240618:word_of_wisdom::5:e2dfcb752d5542d897465de942b36a4d\n
  ```

3. **Client solves the challenge and sends the solution**:
  ```
  <solution>\n
  ```

4. **Server verifies the solution and sends a quote**:
  ```
  The only limit to our realization of tomorrow is our doubts of today.\n
  ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


