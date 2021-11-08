package gts

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex
var t *time.Timer

type Ele struct {
	ID       string
	BootTime int64
	Freq     int64
	Cycles   int64
	Handler  func(interface{})
	Prams    interface{}
}

type Sche interface {
	Add(*Ele) error
	Madd([]*Ele) error
	Mrem([]string) error
	Pop(int) *Ele
	Remove(string) bool
}
type eles []*Ele

func checkEle(new *Ele) error {
	if new.ID == "" {
		return fmt.Errorf("ID empty")
	}
	if new.Freq == 0 {
		return fmt.Errorf("circular time frequency can not be zero")
	}
	if new.Handler == nil {
		return fmt.Errorf("no handle funtion to execute")
	}

	if new.BootTime == 0 {
		new.BootTime = time.Now().Unix() + new.Freq
	}
	return nil
}

func (sche *eles) add(new *Ele) {
	length := len(*sche)
	Idx := length
	pIdx := (Idx+1)/2 - 1
	*sche = append(*sche, new)

	for pIdx >= 0 && (*sche)[Idx].BootTime < (*sche)[pIdx].BootTime {
		(*sche)[pIdx], (*sche)[Idx] = (*sche)[Idx], (*sche)[pIdx]
		Idx = pIdx
		pIdx = (Idx+1)/2 - 1
	}
}

func (sche *eles) pop(index int) *Ele {
	length := len(*sche)
	if length == 0 || index >= length {
		return nil
	}

	res := (*sche)[index]
	(*sche)[index] = (*sche)[length-1]

	for {
		l_idx := index*2 + 1
		r_idx := l_idx + 1

		if l_idx >= length {
			break
		}

		var next int
		if r_idx >= length || (*sche)[l_idx].BootTime < (*sche)[r_idx].BootTime {
			next = l_idx
		} else {
			next = r_idx
		}

		if (*sche)[index].BootTime > (*sche)[next].BootTime {
			(*sche)[index], (*sche)[next] = (*sche)[next], (*sche)[index]
			index = next
		} else {
			break
		}
	}

	*sche = (*sche)[:length-1]
	return res
}

func (sche *eles) resetTimer() {
	if len(*sche) > 0 {
		now := time.Now().Unix()
		nextTime := (*sche)[0].BootTime - now
		if nextTime < 0 {
			nextTime = 0
		}
		t.Reset(time.Duration(nextTime) * time.Second)
	}
}

func (sche *eles) Add(new *Ele) error {
	mutex.Lock()
	defer mutex.Unlock()

	if err := checkEle(new); err != nil {
		return err
	}

	sche.add(new)
	sche.resetTimer()
	return nil
}

func (sche *eles) Pop(index int) *Ele {
	mutex.Lock()
	defer mutex.Unlock()

	task := sche.pop(index)
	sche.resetTimer()
	return task
}

func (sche *eles) Remove(id string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	for i, v := range *sche {
		if v.ID == id {
			sche.pop(i)
			sche.resetTimer()
			return true
		}
	}
	return false
}

func (sche *eles) Madd(news []*Ele) error {
	// check each
	for _, e := range news {
		if err := checkEle(e); err != nil {
			return err
		}
	}

	mutex.Lock()
	defer mutex.Unlock()
	for _, e := range news {
		sche.add(e)
	}
	sche.resetTimer()

	return nil
}

func (sche *eles) Mrem(ids []string) error {
	mutex.Lock()
	defer mutex.Unlock()

	taskMap := make(map[string]int)
	for i, v := range *sche {
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

func ScheInit() Sche {
	elesP := &eles{}
	t = time.NewTimer(10 * time.Second)

	go startSche(elesP)
	return elesP
}

func startSche(sche *eles) {
	for {
		<-t.C
		cur := sche.Pop(0)
		if cur == nil {
			t.Reset(10 * time.Second)
			continue
		}

		go cur.Handler(cur.Prams)

		if cur.Cycles == 1 {
			sche.resetTimer()
			continue
		}
		if cur.Cycles > 1 {
			cur.Cycles--
		}
		cur.BootTime += cur.Freq
		sche.Add(cur)
	}
}
