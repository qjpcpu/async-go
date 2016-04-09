package goasync

import (
	"errors"
	"testing"
)

type MyStruct struct {
	name string
}

func TestAuto(t *testing.T) {
	graph := map[string]*Task{
		"b": &Task{
			Dep: []string{"a"},
			Handler: func(cb Cb, ar ...AsyncResult) {
				t.Log("task b")
				cb("b string", nil)
			},
		},
		"e": &Task{
			Dep: []string{"a", "b", "c"},
			Handler: func(cb Cb, ar ...AsyncResult) {
				var b string
				ar[1].Data(&b)
				if b != "b string" {
					t.Error("should be 'b string'")
				}
				var c MyStruct
				ar[2].Data(&c)
				if c.name != "from c" {
					t.Error("should be 'from c'")
				}
				tbl := map[string]MyStruct{
					"first": MyStruct{name: "inner"},
				}
				cb(tbl, nil)
			},
		},
		"f": &Task{
			Dep: []string{"e"},
			Handler: func(cb Cb, ar ...AsyncResult) {
				var tbl map[string]MyStruct
				ar[0].Data(&tbl)
				if tbl["first"].name != "inner" {
					t.Error("should be 'inner'")
				}
				cb(nil, nil)
			},
		},
		"c": &Task{
			Dep: []string{"a"},
			Handler: func(cb Cb, ar ...AsyncResult) {
				var data []string
				ar[0].Data(&data)
				t.Log("task c get a's data:", data)
				ms := &MyStruct{name: "from c"}
				cb(ms, nil)
			},
		},
		"a": &Task{
			Handler: func(cb Cb, ar ...AsyncResult) {
				t.Log("task a")
				d := []string{"bob", "foo"}
				cb(d, nil)
			},
		},
	}
	asy, _ := Auto(graph)
	asy.Run()
}

func TestParallel(t *testing.T) {
	asy, _ := Parallel(
		func(cb Cb, ar ...AsyncResult) {
			t.Log("aaa")
			cb(0, nil)
		},
		func(cb Cb, ar ...AsyncResult) {
			t.Log("bbb")
			cb("", nil)
		},
	)
	asy.Run()
	names := asy.GetTaskNames()
	var s int = 2
	asy.GetResults(names[1])[0].Data(&s)
	if s != 0 {
		t.Error("should be zero")
	}
	var str string
	asy.GetResults(names[0])[0].Data(&str)
	if str != "" {
		t.Error("should be empty")
	}
}
func TestAutoErr(t *testing.T) {
	graph := map[string]*Task{
		"b": &Task{
			Dep: []string{"a"},
			Handler: func(cb Cb, ar ...AsyncResult) {
				t.Log("task b")
				cb("b string", nil)
			},
		},
		"a": &Task{
			Handler: func(cb Cb, ar ...AsyncResult) {
				t.Log("task a")
				d := []string{"bob", "foo"}
				cb(d, errors.New("error happens in a"))
			},
		},
	}
	asy, _ := Auto(graph)
	err := asy.Run()
	if err == nil {
		t.Error("should get an error")
	}
	arr := asy.GetResults("a")
	if arr[0].err == nil {
		t.Error("should be error")
	}
}
