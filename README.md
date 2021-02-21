# grpc-test

**pingpong_server**

```sh
$ export GRPC_GO_LOG_VERBOSITY_LEVEL=99
$ export GRPC_GO_LOG_SEVERITY_LEVEL=info
$ ./pingpong_server
```

**pingpong_client**

```sh
$ ./pingpong_client -a <server-ip:port> -c 10000
```
