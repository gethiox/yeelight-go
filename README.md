# YeeLight

[![dinosaur with bulb in paw](resources/yee.gif)](https://www.youtube.com/watch?v=q6EoRBvdVPQ)

Go library to control [YEELIGHT](https://www.yeelight.com/) devices for dinosaurs with scary faces

### Disclaimer
- Library is not in stable state and not finished (few missing features)
- User interface may be slightly changed before 1.0 release
- Device auto-discovery are not implemented yet, You'll need to connect to them directly with IP address. 
- Tested only on one type of bulb ([Yeelight Smart LED Bulb (Color)](https://www.yeelight.com/en_US/product/lemon-color)),
  I can't guarantee that everything will work correctly on other devices

# Usage

[![](https://godoc.org/github.com/gethiox/yeelight-go?status.svg)](http://godoc.org/github.com/gethiox/yeelight-go)

make sure LAN control is enabled in your device: 
![](resources/enabling_lan_control.png)

Device initialization:
```go
yl.NewBulb("192.168.0.123")
```

### Available commands

Device functions:
```go
// commands for standard and background light:
func Temperature(temp, duration int) error                                           {}
func RGB(rgb, duration int) error                                                    {}
func HSV(hue, saturation, duration int) error                                        {}
func Brightness(brightness, duration int) error                                      {}
func StartColorFlow(count int, action CfAction, flowExpression FlowExpression) error {}
func StopColorFlow() error                                                           {}
func SetDefault() error                                                              {}
func PowerOn(duration int) error                                                     {}
func PowerOnWithMode(duration int, mode Mode) error                                  {}
func PowerOff(duration int) error                                                    {}
func Toggle() error                                                                  {}

// for background only:
func DevToggle() error {} 

// for standard only:
func CronAdd(jobType CronType, minutes int) error                 {}
func CronDel(jobType CronType) error                              {}
func SetAdjust(action Action, prop AdjustProp) error              {}
func AdjustBright(percentage, duration int) error                 {}
func AdjustTemperature(percentage, duration int) error            {}
func AdjustColor(percentage, duration int) error                  {}
func SetName(name string) error                                   {}
func StartMusic(ifaceName string) (error, musicSupportedCommands) {}

// commands available in music mode
func Temperature(temp, duration int)                                           {}
func RGB(rgb, duration int)                                                    {}
func HSV(hue, saturation, duration int)                                        {}
func Brightness(brightness, duration int)                                      {}
func StartColorFlow(count int, action CfAction, flowExpression FlowExpression) {}
func StopColorFlow()                                                           {}
```

### Example
```go
package main

import ( 
    "fmt"
    "sync"
    "time"

    yl "github.com/gethiox/yeelight-go"
)

func main() {
	var bulbs []*yl.Bulb

	var ipTemplate = "192.168.10.%d"
	var octets = []int{220, 221, 222, 223}

	// Connecting to many bulbs at the same time
	for _, octet := range octets {
		ip := fmt.Sprintf(ipTemplate, octet)
		bulb := yl.NewBulb(ip) 
		err := bulb.Connect()
		if err != nil {
			panic(err)
		}
		bulbs = append(bulbs, bulb)
	}

	jobs := sync.WaitGroup{}
	for _, bulb := range bulbs {
		jobs.Add(1)

		// running goroutine for every bulb separately
		go func(bulb *yl.Bulb) {
			var err error

			err = bulb.PowerOn(0)
			if err != nil {
				panic(err)
			}

			// trying to start music server on "enp0s25" interface
			// empty name ("") may be passed for starting music server on first found interface
			err, music := bulb.StartMusic("enp0s25")
			if err != nil {
				panic(err)
			}

			// go through hue range ten times
			for iterations := 0; iterations < 10; iterations++ {
				for i := 0; i < 360; i++ {
					music.HSV(i, 100, 50)
					time.Sleep(time.Millisecond * 50)
				}
			}

			jobs.Done()
		}(bulb)
	}
	jobs.Wait()
}
```
