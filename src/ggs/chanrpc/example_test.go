package chanrpc_test

import (
	"fmt"
	"ggs/chanrpc"
	"sync"
)

func Example() {
	s := chanrpc.NewServer(10) // 创建RPC服务器

	var wg sync.WaitGroup
	wg.Add(1)

	// goroutine 1
	go func() {
		// 注册id->func的映射
		s.Register("f0", func(args []interface{}) {})
		s.Register("f1", func(args []interface{}) interface{} {
			return 1
		})
		s.Register("fn", func(args []interface{}) []interface{} {
			return []interface{}{1, 2, 3}
		})
		s.Register("add", func(args []interface{}) interface{} {
			n1 := args[0].(int)
			n2 := args[1].(int)
			return n1 + n2
		})

		wg.Done()

		for {
			s.Exec(<-s.ChanCall) // 执行RPC调用
		}
	}()

	wg.Wait()
	wg.Add(1)

	// goroutine 2
	go func() {
		c := s.Open(10) // 创建RPC客户端

		// sync
		err := c.Call0("f0")
		if err != nil {
			fmt.Println(err)
		}

		r1, err := c.Call1("f1")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r1) // 1
		}

		rn, err := c.CallN("fn")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(rn[0], rn[1], rn[2]) // 1, 2, 3
		}

		ra, err := c.Call1("add", 1, 2)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(ra) // 3
		}

		// asyn
		c.AsynCall("f0", func(err error) {
			if err != nil {
				fmt.Println(err)
			}
		})

		c.AsynCall("f1", func(ret interface{}, err error) {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(ret) // 1
			}
		})

		c.AsynCall("fn", func(ret []interface{}, err error) {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(ret[0], ret[1], ret[2]) // 1, 2, 3
			}
		})

		c.AsynCall("add", 1, 2, func(ret interface{}, err error) {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(ret) // 3
			}
		})

		// 执行回调
		c.Cb(<-c.ChanAsynRet)
		c.Cb(<-c.ChanAsynRet)
		c.Cb(<-c.ChanAsynRet)
		c.Cb(<-c.ChanAsynRet)

		// go
		s.Go("f0")

		wg.Done()
	}()

	wg.Wait()

	// Output:
	// 1
	// 1 2 3
	// 3
	// 1
	// 1 2 3
	// 3
}
