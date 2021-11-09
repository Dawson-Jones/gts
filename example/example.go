package main

import (
	"fmt"
	"github.com/Dawson-Jones/gts"
	"log"
	"time"
)

func handler(p interface{}) {
	id, ok := p.(string)
	if !ok {
		return
	}
	fmt.Println(id, time.Now().Unix())
}

func main() {
	var err error
	cron := gts.ScheInit()

	id1 := "111111"
	err = cron.Add(&gts.Ele{
		ID:      id1,
		Freq:    3,
		Handler: handler,
		Prams:   id1,
		Cycles:  2,
	})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id1)

	id2 := "222222"
	err = cron.Add(&gts.Ele{
		ID:      id2,
		Freq:    5,
		Handler: handler,
		Prams:   id2,
		Cycles:  1,
	})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id2)

	id3 := "333333"
	err = cron.Add(&gts.Ele{
		ID:      id3,
		Freq:    7,
		Handler: handler,
		Prams:   id3,
		Cycles:  -1,
	})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id3)

	time.Sleep(20 * time.Second)
	cron.Remove(id1)
	fmt.Println("----------")

	select {}
}

func SaveTODB(id string) {
	// 储存返回的 id
}
