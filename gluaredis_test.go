package gluaredis

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/redis/go-redis/v9"
	lua "github.com/yuin/gopher-lua"
)

func TestNewClient(t *testing.T) {
	if err := evalLua(t, `
        local redis = require("redis")
  
        local rdb = redis.new_client("localhost:6379", "")

        assert_not_equal(nil, rdb)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestSetSuccess(t *testing.T) {
	gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Set", func(_ *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {

		resp := redis.NewStatusCmd(ctx)
		return resp
	})
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rdb.Set(context.Background(), "key", "value", 0)
	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local err = rdb:set("key", "value", 0)
        if err then
            print("set error: ", err)
            return
        end
        assert_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func evalLua(t *testing.T, script string) error {
	L := lua.NewState()
	defer L.Close()

	Preload(L)

	L.SetGlobal("assert_equal", L.NewFunction(func(L *lua.LState) int {
		expected := L.Get(1)
		actual := L.Get(2)

		if expected.Type() != actual.Type() || expected.String() != actual.String() {
			t.Errorf("Expected %s %q, got %s %q", expected.Type(), expected, actual.Type(), actual)
		}

		return 0
	}))

	L.SetGlobal("assert_not_equal", L.NewFunction(func(L *lua.LState) int {
		expected := L.Get(1)
		actual := L.Get(2)

		if expected.Type() == actual.Type() && expected.String() == actual.String() {
			t.Errorf("not expected %s %q, got %s %q", expected.Type(), expected, actual.Type(), actual)
		}

		return 0
	}))

	L.SetGlobal("assert_contains", L.NewFunction(func(L *lua.LState) int {
		contains := L.Get(1)
		actual := L.Get(2)

		if !strings.Contains(actual.String(), contains.String()) {
			t.Errorf("Expected %s %q contains %s %q", actual.Type(), actual, contains.Type(), contains)
		}

		return 0
	}))

	return L.DoString(script)
}
