# gts
multi-tasks schedule  
多任务调度系统  

gts 可以实现多任务的统一管理调度  
动态的`添加 / 删除`定时任务

## 演示
```go
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

```

## 结果
```
dae46767-008a-4862-a872-22a77f074048 1607916994     // 3 秒钟的任务
bb857db5-1718-4457-9671-4de75533d925 1607916996     // 5 秒钟的任务
dae46767-008a-4862-a872-22a77f074048 1607916997     // 3
6e4b230d-cd68-4892-8ae9-3520be53cd9d 1607916998     // 7 秒钟的任务
dae46767-008a-4862-a872-22a77f074048 1607917000     // 3
----------                                          // 删掉了 3 秒钟
bb857db5-1718-4457-9671-4de75533d925 1607917001     // 5
6e4b230d-cd68-4892-8ae9-3520be53cd9d 1607917005     // 7
bb857db5-1718-4457-9671-4de75533d925 1607917006     // 5
bb857db5-1718-4457-9671-4de75533d925 1607917011     // 5
6e4b230d-cd68-4892-8ae9-3520be53cd9d 1607917012     // 7
bb857db5-1718-4457-9671-4de75533d925 1607917016     // 5
6e4b230d-cd68-4892-8ae9-3520be53cd9d 1607917019     // 7

```

## 方法
- Add  
```go
	id1 := uuid.New().String()
	err = cron.Add(&gts.Ele{
		ID:      id1,       // 任务ID
		Freq:    3,         // 频率
		Handler: handler,   // 任务的回调函数
		Prams:   id1,       // 回调函数的参数
		Cycles:  -1,        // 循环次数，负数则是无限循环
	})
	if err != nil {
		log.Fatal(err)
	}
    // 返回
```

- Remove  
return bool  
  *true*: remove successful  
  *false*: remove failed
```go
	cron.Remove(id)     // 
	
	if err != nil {
		log.Fatal(err)
	}
    // 返回
```