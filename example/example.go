package main

import (
	"fmt"
	"github.com/Dawson-Jones/gts"
	"github.com/google/uuid"
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

	id1 := uuid.New().String()
	err = cron.Add(&gts.Ele{
		ID:      id1,
		Freq:    3,
		Handler: handler,
		Prams:   id1,
		Cycles:  -1,
	})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id1)

	id2 := uuid.New().String()
	err = cron.Add(&gts.Ele{
		ID:      id2,
		Freq:    4,
		Handler: handler,
		Prams:   id2,
		Cycles:  -1,
	})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id2)

	id3 := uuid.New().String()
	err = cron.Add(&gts.Ele{
		ID:      id3,
		Freq:    5,
		Handler: handler,
		Prams:   id3,
		Cycles:  -1,
	})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id3)

	time.Sleep(10 * time.Second)
	cron.Remove(id1)
	fmt.Println("----------")

	select {}
}

func SaveTODB(id string) {
	// 储存返回的 id
}
