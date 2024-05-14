package gluaredis

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

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
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(nil)
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local err = rdb:set("key", "value", 0)
        assert_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestSetFail(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(errors.New("test"))
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local err = rdb:set("key", "value", 0)
        assert_not_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestSetNxSuccess(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(nil)
		bcmd, _ := cmd.(*redis.BoolCmd)
		bcmd.SetVal(true)
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local ok, err = rdb:setnx("key", "value", 0)
        assert_equal(err, nil)
        assert_equal(ok, true)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestSetNxFail(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(nil)
		bcmd, _ := cmd.(*redis.BoolCmd)
		bcmd.SetVal(false)
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local ok, err = rdb:setnx("key", "value", 0)
        assert_equal(err, nil)
        assert_equal(ok, false)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestGetSuccess(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(nil)
		scmd, _ := cmd.(*redis.StringCmd)
		scmd.SetVal("value")
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local val, exist, err = rdb:get("key")
        assert_equal(val, "value")
        assert_equal(exist, true)
        assert_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestGetNotExist(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(redis.Nil)
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local val, exist, err = rdb:get("key")
        assert_equal(val, "")
        assert_equal(exist, false)
        assert_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestGetFail(t *testing.T) {
	gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(errors.New("test"))
		return nil
	})
	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local val, exist, err = rdb:get("key")
        assert_equal(val, "")
        assert_equal(exist, false)
        assert_not_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestDelSuccess(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(nil)
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local err = rdb:del("key")
        assert_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestDelFail(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(errors.New("test"))
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local err = rdb:del("key")
        assert_not_equal(err, nil)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestExpireSuccess(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(nil)
		bcmd, _ := cmd.(*redis.BoolCmd)
		bcmd.SetVal(true)
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local ok, err = rdb:expire("key", 0)
        assert_equal(err, nil)
        assert_equal(ok, true)

        `); err != nil {
		t.Errorf("Failed to evaluate script: %s", err.Error())
	}
}

func TestExpireFail(t *testing.T) {
	p := gomonkey.ApplyMethod(reflect.TypeOf(&redis.Client{}), "Process", func(_ *redis.Client, ctx context.Context, cmd redis.Cmder) error {
		cmd.SetErr(nil)
		bcmd, _ := cmd.(*redis.BoolCmd)
		bcmd.SetVal(false)
		return nil
	})
	defer p.Reset()

	if err := evalLua(t, `
        local redis = require("redis")
        
        local rdb = redis.new_client("localhost:6379", "")
        
        
        local ok, err = rdb:expire("key", 0)
        assert_equal(err, nil)
        assert_equal(ok, false)

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
