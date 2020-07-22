package yeelight

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type standardCommands struct {
	commander commander
}

// Prop reads given properties
// Not implemented! TODO: TODO
func (c *standardCommands) Prop(props ...Property) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	return data, errors.New("not implemented")
}

// cronAdd sets timer which invokes given CronType operation (power off is only supported)
func (c *standardCommands) cronAdd(jobType CronType, minutes int) error {
	if jobType != CRON_TYPE_POWER_OFF {
		return errors.New("jobType needs to be 0 (power off/timer)")
	}
	return c.commander.executeCommand(
		partialCommand{"cron_add", params{int(jobType), minutes}},
	)
}

// SetTimer powers off device after given number of minutes
func (c *standardCommands) SetTimer(minutes int) error {
	return c.cronAdd(CRON_TYPE_POWER_OFF, minutes)
}

// Not implemented! TODO: TODO
func (c *standardCommands) CronGet(jobType CronType) error {
	return errors.New("not implemented")
}

// CronDel removes a timer for given CronType operation
func (c *standardCommands) CronDel(jobType CronType) error {
	if !(jobType == CRON_TYPE_POWER_OFF) {
		return errors.New("jobType needs to be 0 (power off/timer)")
	}
	return c.commander.executeCommand(
		partialCommand{"cron_del", params{int(jobType)}},
	)
}

// SetAdjust tunes given AdjustProp in a given Action behavior.
// This method is not very precise, please look for dedicated AdjustXxx functions instead
func (c *standardCommands) SetAdjust(action Action, prop AdjustProp) error {
	if prop == ADJUST_PROP_COLOR && action != ADJUST_ACTION_CIRCLE { // edge case from documentation
		return errors.New("color adjusting can be only performed with \"circle\" action")
	}

	return c.commander.executeCommand(
		partialCommand{"set_adjust", params{string(action), string(prop)}},
	)
}

// AdjustBright adjusts bright, range: -100 - 100
func (c *standardCommands) AdjustBright(percentage, duration time.Duration) error {
	if percentage < -100 || percentage > 100 {
		return errors.New("percentage range must be -100 - 100")
	}

	return c.commander.executeCommand(
		partialCommand{"adjust_bright", params{percentage, timeToMs(duration)}},
	)
}

// AdjustTemperature adjusts temperature, range: -100 - 100
func (c *standardCommands) AdjustTemperature(percentage, duration time.Duration) error {
	if percentage < -100 || percentage > 100 {
		return errors.New("percentage range must be -100 - 100")
	}

	return c.commander.executeCommand(
		partialCommand{"adjust_ct", params{percentage, timeToMs(duration)}},
	)
}

// AdjustColor adjusts color, range: -100 - 100
func (c *standardCommands) AdjustColor(percentage, duration time.Duration) error {
	if percentage < -100 || percentage > 100 {
		return errors.New("percentage range must be -100 - 100")
	}

	return c.commander.executeCommand(
		partialCommand{"adjust_color", params{percentage, timeToMs(duration)}},
	)
}

// SetName sets device name
func (c *standardCommands) SetName(name string) error {
	return c.commander.executeCommand(
		partialCommand{"set_name", params{name}},
	)
}

func findIface(name string) ([]net.Interface, error) {
	var ifacesToReturn []net.Interface

	interfaces, err := net.Interfaces()
	if err != nil {
		return []net.Interface{}, fmt.Errorf("iface not found: %v", err)
	}

	for _, iface := range interfaces {
		if iface.Name == name {
			ifacesToReturn = append(ifacesToReturn, iface)
		}
	}
	if len(ifacesToReturn) == 0 {
		return ifacesToReturn, fmt.Errorf("iface \"%s\" not found", name)
	}
	return ifacesToReturn, nil
}

// findIPv4Addr finds first valid IPv4 address on given net.Interface
func findIPv4Addr(iface net.Interface) (net.IP, error) {
	var retAddr net.IP

	addresses, err := iface.Addrs()
	if err != nil {
		return retAddr, fmt.Errorf("failed to fetch binded addresses on \"%s\" interface: %v", iface.Name, err)
	}

	for _, addr := range addresses {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}

		// check if ip is IPv4 type
		if ip.To4() != nil {
			return ip, nil
		}
	}
	return net.IP{}, errors.New(fmt.Sprintf("IPv4 not found on \"%s\" interfgace", iface.Name))
}

func openSocket(ip net.IP) (net.Listener, error) {
	// var ipAddr = "" // binding on all interfaces
	address := fmt.Sprintf("%s:0", ip.String()) // :0 for automatic port selection

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// StartMusic starts tries to run music mode.
// You can perform operations on returned music object without quota limitations
// Interface name can be passed to select exact interface for music server on first assigned IPv4 address
// (bulb needs to connect to opened socket by client), empty string may be passed ("") for
// trying to connect on first available (up and non-loopback) interface and first assigned IPv4 address
func (c *standardCommands) StartMusic() (musicSupportedCommands, error) {
	var (
		ifacesToTry []net.Interface
		err         error
	)

	ifacesToTry, err = net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to read available interfaces: %v", err)
	}

	var (
		binded           bool
		bindedIface      net.Interface
		bindedIPv4Addr   net.IP
		bindedConnection *net.TCPAddr
		listener         net.Listener
	)

	for _, iface := range ifacesToTry {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // Interface is neither up or non-loopback (localhost)
		}

		bindedIPv4Addr, err = findIPv4Addr(iface)
		if err != nil {
			continue // failing find a valid IPv4 address
		}

		// TODO: subnet validation could be invoked here
		listener, err = openSocket(bindedIPv4Addr) // first 1024 ports are root-only
		if err != nil {
			continue // port opening failed
		}

		tcpAddr, ok := listener.Addr().(*net.TCPAddr)
		if !ok {
			return nil, errors.New("listener address is somehow not a *net.TCPAddr type")
		}

		bindedConnection = tcpAddr
		bindedIface = iface
		binded = true
		break
	}

	if !binded {
		return nil, fmt.Errorf("failed to bind on any of given interfaces")
	}

	log.Printf("[music] binded on \"%v\" iface on \"%s\" address on \"%d\" port",
		bindedIface.Name, bindedIPv4Addr, bindedConnection.Port)

	incomingConnection := make(chan net.Conn)
	defer close(incomingConnection)

	// starting "music server" and waiting for first incoming connection
	go func(listener net.Listener) {
		log.Printf("[music] Waiting for a device connection...")
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[music] Device connection failed: %v", err)
			return
		}
		log.Printf("[music] Device connected!")
		incomingConnection <- conn

		var buf = make([]byte, 1)
		go func(conn net.Conn) { // not very clever solution but at least something for debugging
			_, err = conn.Read(buf)
			log.Printf("[music] Device disconnected")
		}(conn)
	}(listener)

	log.Printf("[music] Initializating Music Mode...")
	err = c.commander.executeCommand(
		partialCommand{"set_music", params{1, bindedIPv4Addr, bindedConnection.Port}},
	)
	// err main contains "map[code:-5001 message:invalid params]" if Music mode is already running
	if err != nil {
		return nil, err
	}

	select {
	case conn := <-incomingConnection:
		music := NewMusic(conn)
		if music == nil {
			return nil, errors.New("[music] Connection failed")
		}
		log.Printf("[music] Music Mode Initialized!")
		return music, nil
	case <-time.After(time.Second * 2):
		err := listener.Close()
		if err != nil {
			log.Printf("[music] failed to close music server: %v", err)
		}
		return nil, fmt.Errorf("device connection timeout")
	}
}

// chooseEffect returns effect string Accordingly to given duration value.
// Chooses between "sudden" and "smooth".
func chooseEffect(duration time.Duration) (string, error) {
	// if duration != 0 && duration < 30 {
	//	return "", errors.New("Ooops")
	// }

	if duration == time.Duration(0) {
		return "sudden", nil
	}
	return "smooth", nil
}
