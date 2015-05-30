package physical

import (
    "fmt"
    "time"

    "github.com/armon/go-metrics"
    "github.com/keimoon/gore"
)

// A backend to store key/value pairs in redis
type RedisBackend struct {
    // path is used as a key prefix
    path string
    client *gore.Conn
}

// newRedisBackend constructs a new backend using the given server address
func newRedisBackend(conf map[string]string) (Backend, error) {

    // Get or set path. Defaults to vault:
    path, ok := conf["path"]
    if !ok {
        path = "vault:"
    }

    // Get or set reddis address. Defaults to the localhost and default port
    address, ok := conf["address"]
    if !ok {
        address = "127.0.0.1:6379"
    }

    redisConn, err := gore.Dial(address)
    if err != nil {
        fmt.Errorf("Unable to connect to redis server at '%s'", address)
    }

    r := &RedisBackend {
        client: redisConn,
        path: path,
    }
    return r, nil
}

// Put is used to insert or update an entry
func (r *RedisBackend) Put(entry *Entry) error {
    defer metrics.MeasuredSince([]string{"redis", "put"}, time.Now())

    _, err := gore.NewCommand("SET", entry.Key, entry.Value).Run(r.client)
    if err != nil {
        fmt.Errorf("Unable to create entry: ", err)
    }
    return nil
}

// TODO implement Get

// Delete is used to destroy an entry in the backend
func (r *RedisBackend) Delete(key string) error {
    _, err := gore.NewCommand("DEL", key).Run(r.client)
   if err != nil {
       fmt.Errorf("Unable to delete entry: ", err)
   }
   return nil
}

// TODO implement List
