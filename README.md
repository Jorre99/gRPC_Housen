# gRPC_Housen
Grpc Server-clients

## goclient

### Install

`$ go install github.com/Jorre99/gRPC_Housen/carabiner`

### Use

```
$ export CHATTER_USER=${user}  # Your own username.
$ export CHATTER_FRIEND=${friend}  # Person you want to chat with.
$ carabiner 
```

You can use the `CHATTER_SERVER` environment variable to select a server.

## JavaClient

###Generate Proto

````
$ protoc -I . --java_out=server_fllower_house server_fllower_house/proto/HouseServer.proto

````

## Python Server

###Generate Proto

````
$ python3 -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. server_fllower_house/proto/HouseServer.proto
````
