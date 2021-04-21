# go-list

## Introduction

go-list是一个支持并发操作的有序链表.

## Features

- 并发安全.
- 无锁查询和遍历.
- 元素有序.

## QuickStart

```go
package main

import (
	"fmt"

	"github.com/ytz12345/golist"
)

func main() {
	l := golist.NewInt()

	for _, v := range []int{10, 12, 15} {
		if l.Insert(v) {
			fmt.Println("golist insert", v)
		}
	}

	if l.Contains(10) {
		fmt.Println("golist contains 10")
	}

	l.Range(func(value int) bool {
		fmt.Println("golist range found ", value)
		return true
	})

	l.Delete(15)
	fmt.Printf("golist contains %d items\r\n", l.Len())
}

```