package gts

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

var mutex sync.Mutex
var t *time.Timer

type Ele struct {
	ID       string
	BootTime int64
	Freq     int64
}

type Eles []*Ele

type Cron struct {
	C     chan *Ele
	Tasks Eles
}

func checkEle(new *Ele) error {
	if new.ID == "" {
		new.ID = uuid.New().String()
	}
	if new.Freq == 0 {
		return fmt.Errorf("circular time frequency can not be zero")
	}
	if new.BootTime == 0 {
		new.BootTime = time.Now().Unix() + new.Freq
	}

	return nil
}

func (sche *Cron) add(new *Ele) string  {

	length := len(sche.Tasks)
	Idx := length
	pIdx := (Idx+1)/2 - 1
	sche.Tasks = append(sche.Tasks, new)

	for ; pIdx >= 0 && sche.Tasks[Idx].BootTime < sche.Tasks[pIdx].BootTime; {
		sche.Tasks[pIdx], sche.Tasks[Idx] = sche.Tasks[Idx], sche.Tasks[pIdx]
		Idx = pIdx
		pIdx = (Idx+1)/2 - 1
	}

	return new.ID
}

func (sche *Cron) pop(index int) *Ele {
	length := len(sche.Tasks)
	if length == 0 || index >= length {
		return nil
	}

	res := sche.Tasks[index]
	sche.Tasks[index] = sche.Tasks[length-1]

	for {
		l_idx := index*2 + 1
		r_idx := l_idx + 1

		if l_idx >= length {
			break
		}

		var next int
		if r_idx >= length || sche.Tasks[l_idx].BootTime < sche.Tasks[r_idx].BootTime {
			next = l_idx
		} else {
			next = r_idx
		}

		if sche.Tasks[index].BootTime > sche.Tasks[next].BootTime {
			sche.Tasks[index], sche.Tasks[next] = sche.Tasks[next], sche.Tasks[index]
			index = next
		} else {
			break
		}
	}

	sche.Tasks = sche.Tasks[:length-1]
	return res
}

func (sche *Cron) resetTimer() {
	if len(sche.Tasks) > 0 {
		now := time.Now().Unix()
		nextTime := sche.Tasks[0].BootTime - now
		if nextTime < 0 {
			nextTime = 0
		}
		t.Reset(time.Duration(nextTime) * time.Second)
	}
}

func (sche *Cron) Add(new *Ele) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if err := checkEle(new); err != nil {
		return "", err
	}

	scheID := sche.add(new)
	sche.resetTimer()
	return scheID, nil
}

func (sche *Cron) Pop(index int) *Ele {
	mutex.Lock()
	defer mutex.Unlock()

	task := sche.pop(index)
	sche.resetTimer()
	return task
}

func (sche *Cron) Remove(id string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	for i, v := range sche.Tasks {
		if v.ID == id {
			sche.pop(i)
			sche.resetTimer()
			return true
		}
	}
	return false
}

func (sche *Cron) Madd(news Eles) ([]string, error) {
	// check each
	for _, e := range news {
		if err := checkEle(e); err != nil {
			return nil, err
		}
	}

	mutex.Lock()
	defer mutex.Unlock()
	var res []string
	for _, e := range news {
		res = append(res, sche.add(e))
	}
	sche.resetTimer()

	return res, nil
}

func (sche *Cron) Mrem(ids []string) error {
	mutex.Lock()
	defer mutex.Unlock()

	taskMap := make(map[string]int)
	for i, v := range sche.Tasks {
		taskMap[v.ID] = i
	}

	for _, id := range ids {
		if _, ok := taskMap[id]; !ok {
			return fmt.Errorf("%s not in Crontabs", id)
		}
	}

	for _, id := range ids {
		sche.pop(taskMap[id])
	}
	sche.resetTimer()
	return nil
}

func NewCron() *Cron {
	ch := make(chan *Ele, 5)
	sche := &Cron{
		C:     ch,
		Tasks: Eles{},
	}

	go startCron(sche)
	return sche
}

func startCron(sche *Cron) {
	t = time.NewTimer(10 * time.Second)

	for {
		select {
		case <-t.C:
			cur := sche.Pop(0)
			if cur == nil {
				t.Reset(10 * time.Second)
				continue
			}

			cur.BootTime += cur.Freq
			sche.Add(cur)
			sche.C <- cur
		}
	}
}
