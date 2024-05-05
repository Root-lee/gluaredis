package gluaredis

import (
	"github.com/redis/go-redis/v9"

	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("redis", Loader)
}

var exports = map[string]lua.LGFunction{
	"new_client": newClient,
}

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)

	registerRedisClientType(L)

	return 1
}

func newClient(L *lua.LState) int {
	addr := L.CheckString(1)
	password := L.CheckString(2)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	ud := L.NewUserData()
	ud.Value = rdb
	L.SetMetatable(ud, L.GetTypeMetatable(REDIS_CLIENT_TYPENAME))
	L.Push(ud)
	return 1
}
