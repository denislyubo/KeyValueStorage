# KeyValueStorage

## Add pair {any_key, any_value}
```curl -X PUT -d "any_value" http://localhost:8080/v1/any_key```

## Get key {any_key}
```curl http://localhost:8080/v1/any_key```

## Add pair {any_key, any_value}
```go run kvs_client.go put any_key any_value```

## Get key {any_key}
```go run kvs_client.go get any_key```
