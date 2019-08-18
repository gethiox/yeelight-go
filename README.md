# Yee Light

![](resources/yee.gif)


# Usage

Device auto-discovery are not implemented yet, You'll need to connect to them directly with IP address. 

Example
```go
package main

import ( 
    "fmt"
    "sync"
    "time"

    yl "github.com/gethiox/go-yee-light-go"
)

func main() {
	var ipTemplate = "192.168.10.%d"
	var octets = []int{220, 221, 222, 223}
	var bulbs []*yl.Bulb

	for _, octet := range octets {
		ip := fmt.Sprintf(ipTemplate, octet)
		bulb := yl.NewBulb(ip, 55443)
		err := bulb.Connect()
		if err != nil {
			panic(err)
		}
		bulbs = append(bulbs, bulb)
	}

	gorutines := sync.WaitGroup{}
	for _, bulb := range bulbs {
		gorutines.Add(1)

		go func(bulb *yl.Bulb) {
			var err error

			err = bulb.PowerOn(0)
			if err != nil {
				panic(err)
			}

			err, music := bulb.StartMusic()
			if err != nil {
				panic(err)
			}

			for iterations := 0; iterations < 10; iterations++ {
				for i := 0; i < 360; i++ {
					music.HSV(i, 100, 50)
					time.Sleep(time.Millisecond * 50)
				}
			}

			gorutines.Done()
		}(bulb)
	}
	gorutines.Wait()
}
```