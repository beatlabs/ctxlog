# Context Logger

`ctxlog` is a library which helps to add a contextual data to your log messages at any time, and have it logged with each message.

## Short-term plans 

_TBD but be aware it might happen_

- Remove dependency on Patron and replace it with some generic logger
- Add examples in `examples/` folder

## Install it

To add this library as a dependency of your project, run 
```shell
$ go get github.com/beatlabs/ctxlog
```

## Usage examples

Initiate a logger in a request context

```go
req := &http.Request{}

ctx := ctxlog.AddLoggerForRequest(req)
```

Add some custom data to it
```go
ctxlog.FromContext(ctx).Int("answer", 42)
ctxlog.FromContext(ctx).Str("so", "long")
ctxlog.FromContext(ctx).SubCtx(map[string]interface{}{
    "thanks_for":  "all the fish",
    "planet":      "Earth",
    "happened_on": "Thursday",
})
```

And log it
```go
ctxlog.FromContext(ctx).Warnf("this %s has — or rather had — a %s, which was this: %s", "planet", "problem",
	"most of the people living on it were unhappy for pretty much all of the time")
```
