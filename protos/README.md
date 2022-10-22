# Proto definitions

Protocol buffers

## Go SDK gen

To generate the Go SDK run:

```bash
protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./tb.proto
```

## JS SDK gen

```bash
protoc -I=./ geo_service.proto --js_out=import_style=commonjs:./
```
