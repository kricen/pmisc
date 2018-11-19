package main

import (
	"fmt"
	"time"

	"github.com/pmisc/prometheus/customized"
	"github.com/pmisc/registry"
)

// this main package aim at to help you how to use it

func main() {

	cc := customized.NewCpuCollector()
	dc := customized.NewDiskCollector()
	mc := customized.NewMemoryCollector()
	rc := customized.NewRequestCollector()

	cr, err := registry.NewCollectorRegister("monitor-helper-test", []string{"http://hello:2379", "http://hello2:2379"})
	if err != nil {
		fmt.Println(err.Error())
	}
	cr.Registe(cc)
	cr.Registe(dc)
	cr.Registe(mc)
	cr.Registe(rc)
	fmt.Println(cr.ToString())
	cr.Start()
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(100 * time.Millisecond)
			rc.AddRecord("getUserName", int64(i))
		}
	}()

	for {
		time.Sleep(10 * time.Second)
		// cr.Push()
	}

	// cli, err := clientv3.New(clientv3.Config{
	// 	Endpoints:   []string{"http://hell03:2379"},
	// 	DialTimeout: 5 * time.Second,
	// })
	// if err != nil {
	// 	// handle error!
	// 	fmt.Println(err.Error())
	// 	return
	// }
	//
	// ctx := context.TODO()
	// ch := cli.Watch(ctx, "/key/", clientv3.WithPrefix())
	// for {
	// 	log.Print("rev")
	// 	select {
	// 	case c := <-ch:
	// 		for _, e := range c.Events {
	// 			log.Printf("%+v", e)
	// 		}
	// 	}
	// }
}
