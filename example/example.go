package main

import (
	"fmt"
	"github.com/Dawson-Jones/gts"
	"log"
	"time"
)

func main() {
	cron := gts.NewCron()

	id1, err := cron.Add(&gts.Ele{Freq: 3})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id1)

	id2, err := cron.Add(&gts.Ele{Freq: 5})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id2)

	id3, err := cron.Add(&gts.Ele{Freq: 7})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id3)

	go func() {
		for {
			select {
			case t := <-cron.C:
				fmt.Println(t.ID, time.Now().Unix())
			}
		}
	}()

    time.Sleep(10 * time.Second)
	cron.Remove(id1)
	fmt.Println("----------")

	select {}
}

func SaveTODB(id string) {
	// 储存返回的 id
}