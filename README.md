# gRPC communication patterns

This project demonstrates various gRPC communication patterns in Go service:

* Unary RPC
* Server Streaming
* Client Streaming
* Bidirectional Streaming

---

## Proto Definition

File: `main.proto`

```proto
syntax = "proto3";

package deepdive;

option go_package = "/proto/gen;mainpb";

service Deepdive{
    rpc Add(AddRequest) returns(AddResponse);
    rpc GenerateFibonacci (FibonacciRequest) returns(stream FibonacciResponse); // server stream
    rpc SendNumbers(stream NumberRequest) returns(NumberResponse); // client stream
    rpc Chat(stream ChatMessage) returns(stream ChatMessage); // bidirectional stream
}

message NumberRequest{
    int32 number = 1;
}

message NumberResponse{
    int32 sum = 2;
}
message AddRequest{
    int32 a = 1;
    int32 b = 2;
}

message AddResponse{
    int32 sum=1;
}

message FibonacciRequest{
    int32 n = 1;
}

message FibonacciResponse{
    int32 number = 1;
}

message ChatMessage{
    string message = 1;
}
```

---

## Running the Server

Make sure your gRPC server is running and listening on port `50051`.

```bash
go mod tidy
go run main.go
```

Expected log:

```bash
Listening on port : 50051
```

---

## Using `grpcurl` to Test

Ensure `grpcurl` is installed. You can install it via:

```bash
sudo snap install grpcurl --edge
```

---

### Case 1: Unary – Add Two Numbers

**Request:**

```bash
grpcurl -plaintext -d '{"a": 3, "b": 4}' localhost:50051 deepdive.Deepdive/Add
```

**Response:**

```json
{
  "sum": 7
}
```

---

### Case 2: Server Streaming – Generate Fibonacci Numbers

**Request:**

```bash
grpcurl -plaintext -d '{"n":10}' localhost:50051 deepdive.Deepdive/GenerateFibonacci
```

**Response:**

```json
{
  "number": 1
}
{
  "number": 1
}
{
  "number": 2
}
{
  "number": 3
}
{
  "number": 5
}
{
  "number": 8
}
{
  "number": 13
}
{
  "number": 21
}
{
  "number": 34
}
```

---

### Case 3: Bidirectional Streaming – Chat

**Request:**

```bash
grpcurl -plaintext -d @ localhost:50051 deepdive.Deepdive/Chat
```

You’ll enter interactive mode. Type each message and hit Enter:

```
{"message":"hello"}
```

**Expected Server Log:**

```bash
 2025/10/02 11:35:57 Received from client: hello
```

You can send message from `server` to `client`:

```bash
{"message":"hi"}
```

---

### Case 4: Client Streaming – Send Multiple Numbers

**Request:**

```bash
grpcurl -plaintext -d @ localhost:50051 deepdive.Deepdive/SendNumbers
```

Then enter the following messages one per line:

```
{"number":1}
{"number":2}
{"number":3}
{"number":4}
{"number":5}
```

Finish the input with `Ctrl+D` to send EOF and trigger the server response.

**Response:**

```json
{
  "sum": 15
}
```

## Working

https://github.com/user-attachments/assets/7a50d203-73be-4834-9235-332532366d98

