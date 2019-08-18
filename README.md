# Yee Light

[![dinosaur with bulb in paw](resources/yee.gif)](https://www.youtube.com/watch?v=q6EoRBvdVPQ)


# Usage

Device auto-discovery are not implemented yet, You'll need to connect to them directly with IP address. 

Example
```go
package main

import ( 
    "fmt"
    "sync"
    "time"

    yl "github.com/gethiox/yeelight-go"
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

			err, music := bulb.StartMusic("192.168.10.100") // your (client's) ip address
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

### Available commands

Bulb functions:
```go
// commands for standard and background light:
func Temperature(temp, duration int) error 
func RGB(rgb, duration int) error 
func HSV(hue, saturation, duration int) error 
func Brightness(brightness, duration int) error 
func StartColorFlow(count int, action CfAction, flowExpression FlowExpression) error 
func StopColorFlow() error 
func SetScene(scene Scene) error // not implemented
func SetDefault() error 
func PowerOn(duration int) error 
func PowerOnWithMode(duration int, mode Mode) error 
func PowerOff(duration int) error 
func Toggle() error

// background only:
func DevToggle() error 

// standard only:
func Prop(props ...Property) (map[string]interface{}, error) // not implemented
func CronAdd(jobType CronType, minutes int) error 
func CronGet(jobType CronType) error // not implemented
func CronDel(jobType CronType) error 
func SetAdjust(action Action, prop AdjustProp) error 
func AdjustBright(percentage, duration int) error 
func AdjustTemperature(percentage, duration int) error 
func AdjustColor(percentage, duration int) error 
func SetName(name string) error 
func StartMusic(hostIP string) (error, musicSupportedCommands)
```

Music functions: (performed on object returned by `StartMusic()`)
```go
func Temperature(temp, duration int)
func RGB(rgb, duration int)
func HSV(hue, saturation, duration int)
func Brightness(brightness, duration int)
func StartColorFlow(count int, action CfAction, flowExpression FlowExpression)
func StopColorFlow()
```

More detailed documentation is available in source code - should be extracted and available in IDE of your choice.