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

	"github.com/ytz123/golist"
)

func main() {
	l := NewInt()

	for _, v := range []int{10, 12, 15} {
		if l.Add(v) {
			fmt.Println("skipset add", v)
		}
	}

	if l.Contains(10) {
		fmt.Println("skipset contains 10")
	}

	l.Range(func(value int) bool {
		fmt.Println("skipset range found ", value)
		return true
	})

	l.Remove(15)
	fmt.Printf("skipset contains %d items\r\n", l.Len())
}

```