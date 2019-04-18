
### The signal server of [hlsjs-p2p-engine](https://github.com/cdnbye/hlsjs-p2p-engine)

This is a distributed signaling implemented using RPC service,you can deploy the code to multiple servers，so that it can handle more signaling.

### 1.install the redis service
### 2.modify the redisconst.go,update the redis host info
### 3.of course,you can change others configuration :
- 3.1 the default server port is 8082,you can modify the serverconst.go and then change the server port;

- 3.2 the default rpc service post is 9002,you can modify the rpcconst.go and the change the rpc service port;

- 3.2 you can modify the logconst.go and the change the log configuration;

### 4.build
- 4.1 go to the main directory of the code

- 4.2 go build main.go

- 4.2 go run main.go

###[hlsjs-p2p-engine](https://github.com/cdnbye/hlsjs-p2p-engine)
