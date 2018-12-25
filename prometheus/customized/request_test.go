package customized

import (
	"fmt"
	"testing"
)

func TestDataRace(t *testing.T) {
	// c := NewRequestCollector()
	c := NewCpuCollector()
	// go func() {
	// 	for {
	// 		c.AddRecord("hello", 100)
	// 		fmt.Println("-------")
	// 	}
	//
	// }()
	for {
		// time.Sleep(10 * time.Millisecond)
		c.Collect()
		fmt.Println("*******")
	}

}
