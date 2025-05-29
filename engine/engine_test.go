package engine_test

import (
	"skii-db/engine"
	"testing"
)

func Test_SetGetKeyValue(t *testing.T) {
    e, _ := engine.NewEngine()
    e.Set("foo", "bar")
    value, err := e.Get("foo")
    if err != nil {
        t.Error(err)
    }
    if value != "bar" {
        t.Error("value should be bar")
    }

    _, err = e.Get("notfound")
    if err == nil {
        t.Error("should return error")
    }

}