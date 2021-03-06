package main

import (
	"context"
	"fmt"
	"time"
	"unsafe"
)

func main() {
	// build a context for test
	ctx := context.Background()

	ctx = context.WithValue(ctx, "key1", "value1")

	ctx, _ = context.WithCancel(ctx)

	ctx = context.WithValue(ctx, &TestData{
		Name:   "key2",
		Number: 1,
	}, "value2")

	ctx, _ = context.WithDeadline(ctx, time.Now().Add(time.Second*20))

	ctx = context.WithValue(ctx, &TestData{
		Name:   "key3",
		Number: 2,
	}, "value3")

	ctx = context.WithValue(ctx, &TestData{
		Name:   "key4",
		Number: 3,
	}, nil)

	ctx = context.WithValue(ctx, &TestData{
		Name:   "key5",
		Number: 4,
	}, nil)

	// get keys and values
	m := GetKeyValues(ctx)
	printMapKeyValue(m)

	// output:
	// [key: &{Name:key3 Number:2}] [value: value3]
	// [key: &{Name:key2 Number:1}] [value: value2]
	// [key: key1] [value: value1]
	// [key: &{Name:key5 Number:4}] [value: <nil>]
	// [key: &{Name:key4 Number:3}] [value: <nil>]
}

type TestData struct {
	Name   string
	Number int
}

type iface struct {
	itab, data uintptr
}

type emptyCtx int

type valueCtx struct {
	context.Context
	key, val interface{}
}

func GetKeyValues(ctx context.Context) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	getKeyValue(ctx, m)
	return m
}

func getKeyValue(ctx context.Context, m map[interface{}]interface{}) {
	ictx := *(*iface)(unsafe.Pointer(&ctx))
	if ictx.data == 0 || int(*(*emptyCtx)(unsafe.Pointer(ictx.data))) == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil {
		m[valCtx.key] = valCtx.val
	}
	getKeyValue(valCtx.Context, m)
}

func printMapKeyValue(m map[interface{}]interface{}) {
	for k, v := range m {
		fmt.Printf("[key: %+v] [value: %+v]\n", k, v)
	}
}
