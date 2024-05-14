# gluaredis
a redis client for [GopherLua](https://github.com/yuin/gopher-lua) VM, based on [redis/go-redis](https://github.com/redis/go-redis ) 

## Installation 
```bash 
go get github.com/Root-lee/gluaredis
``` 

## Using
### Loading Modules 
```go
package main

import (
	"github.com/Root-lee/gluaredis"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	gluaredis.Preload(L)
	defer L.Close()

	if err := L.DoFile("test.lua"); err != nil {
		panic(err)
	}
}
```

### Using In lua <a name="lua-demo-anchor"></a>

test.lua
```lua
local redis = require("redis")

local rdb = redis.new_client("localhost:6379", "")


local err = rdb:set("name", "Root-lee", 0)
print("set:", err)

local success, err = rdb:setnx("name2", "Root-lee2", 0)
print("setnx:", success, err)

local val, exist, err = rdb:get("name")
print("get:", val, exist, err)

local success, err = rdb:expire("name", 1)
print("expire:", success, err)
local val, exist, err = rdb:get("name")
print("get:", val, exist, err)

local err = rdb:del("name")
print("del:", err)
local val, exist, err = rdb:get("name")
print("get:", val, exist, err)
```

## Testing

### Unit Test
```bash
$go test -gcflags="all=-N -l" github.com/Root-lee/gluaredis...
ok      github.com/Root-lee/gluaredis   0.342s
```

### Manual Test
You can start a local redis server: localhost:6379

You can refer to this [lua script](#lua-demo-anchor) to write your own lua script

## License

MIT
