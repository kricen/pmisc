package lib

import (
	"log"
	"testing"
)

func TestJob(t *testing.T) {
	jq := NewJobQueue(5)
	for i := 0; i < 6; i++ {
		p := jq.Push(i)
		if p != nil {
			log.Println("---", p)
		}
	}
	for _, v := range jq.datas {
		log.Println(v)
	}

}
