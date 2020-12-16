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
	"log"
	"time"
)

func main() {
	cron := gts.NewCron()

	id, err := cron.Add(&gts.Ele{Freq: 3})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id)

	id, err = cron.Add(&gts.Ele{Freq: 5})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id)

	id, err = cron.Add(&gts.Ele{Freq: 7})
	if err != nil {
		log.Fatal(err)
	}
	SaveTODB(id)

	go func() {
		for {
			select {
			case t := <-cron.C:
				fmt.Println(t.ID, time.Now().Unix())
			}
		}
	}()
	
    time.Sleep(10 * time.Second)
	cron.Remove(id)
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
	cron := gts.NewCron()
	
	id, err := cron.Add(&gts.Ele{
        ID       string     // 可选项， 不填则自动生成
        BootTime int64      // 可选， 默认为当前时间的 Freq 秒之后
		Freq: 3             // 必填项, 循环触发的频率
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