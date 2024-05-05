package gluaredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	lua "github.com/yuin/gopher-lua"
)

const (
	REDIS_CLIENT_TYPENAME = "redis_client_typename"
)

var ctx = context.Background()

func registerRedisClientType(L *lua.LState) {
	mt := L.NewTypeMetatable(REDIS_CLIENT_TYPENAME)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), redisClientMethods))
}

var redisClientMethods = map[string]lua.LGFunction{
	"set":    set,
	"setnx":  setnx,
	"get":    get,
	"del":    del,
	"expire": expire,
}

func checkRedisClient(L *lua.LState) *redis.Client {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*redis.Client); ok {
		return v
	}
	L.ArgError(1, "redis client expected")
	return nil
}

func set(L *lua.LState) int {
	rdb := checkRedisClient(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	expire := L.CheckInt(4)
	err := rdb.Set(ctx, key, value, time.Duration(expire)*time.Second).Err()

	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

func setnx(L *lua.LState) int {
	rdb := checkRedisClient(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	expire := L.CheckInt(4)
	success, err := rdb.SetNX(ctx, key, value, time.Duration(expire)*time.Second).Result()

	L.Push(lua.LBool(success))
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 2
}

func get(L *lua.LState) int {
	rdb := checkRedisClient(L)
	key := L.CheckString(2)

	val, err := rdb.Get(ctx, key).Result()
	L.Push(lua.LString(val))

	if err == redis.Nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LNil)
	} else if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LBool(true))
		L.Push(lua.LNil)
	}
	return 3
}

func del(L *lua.LState) int {
	rdb := checkRedisClient(L)
	key := L.CheckString(2)
	err := rdb.Del(ctx, key).Err()

	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

func expire(L *lua.LState) int {
	rdb := checkRedisClient(L)
	key := L.CheckString(2)
	expire := L.CheckInt(3)
	success, err := rdb.Expire(ctx, key, time.Duration(expire)*time.Second).Result()

	L.Push(lua.LBool(success))
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 2
}
