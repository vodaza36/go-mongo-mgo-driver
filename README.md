# go-mongo-mgo-driver

In this example the basic usage of the MGO Mongo driver will be shown.

## Documents

- MGO Driver: <https://labix.org/mgo>
- API Docs: <https://godoc.org/gopkg.in/mgo.v2>

## Features

- Insert records
- Find record
- Ensure Index
- Parallel inserts (go routines)
- Tracing

## Tracing

Start the trace analyze tool:

```go
go tool trace trace.out
```