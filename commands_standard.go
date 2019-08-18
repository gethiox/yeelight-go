package yeelight

import (
	"errors"
	"log"
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

// CronAdd sets timer which invokes given CronType operation (power off is only supported)
func (c *standardCommands) CronAdd(jobType CronType, minutes int) error {
	if !(jobType == CRON_TYPE_POWER_OFF) {
		return errors.New("jobType needs to be 0 (power off/timer)")
	}
	return c.commander.executeCommand(
		partialCommand{"cron_add", params{int(jobType), minutes}},
	)
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
func (c *standardCommands) AdjustBright(percentage, duration int) error {
	if percentage < -100 || percentage > 100 {
		return errors.New("percentage range must be -100 - 100")
	}

	return c.commander.executeCommand(
		partialCommand{"adjust_bright", params{percentage, duration}},
	)
}

// AdjustTemperature adjusts temperature, range: -100 - 100
func (c *standardCommands) AdjustTemperature(percentage, duration int) error {
	if percentage < -100 || percentage > 100 {
		return errors.New("percentage range must be -100 - 100")
	}

	return c.commander.executeCommand(
		partialCommand{"adjust_ct", params{percentage, duration}},
	)
}

// AdjustColor adjusts color, range: -100 - 100
func (c *standardCommands) AdjustColor(percentage, duration int) error {
	if percentage < -100 || percentage > 100 {
		return errors.New("percentage range must be -100 - 100")
	}

	return c.commander.executeCommand(
		partialCommand{"adjust_color", params{percentage, duration}},
	)
}

// SetName sets device name
func (c *standardCommands) SetName(name string) error {
	return c.commander.executeCommand(
		partialCommand{"set_name", params{name}},
	)
}

// StartMusic starts music mode. hostIP (client IP) is required to tell device where to connect.
// You can perform operations on returned music object without quota limitations
func (c *standardCommands) StartMusic(hostIP string) (error, musicSupportedCommands) {
	listener, port, err := openSocket("", 1023, 1<<16-1) // first 1024 ports are root-only
	if err != nil {
		return err, nil
	}

	incoming := make(chan *Music)
	defer close(incoming)

	go func() {
		log.Printf("[music] Waiting for a device connection...")
		conn, err := listener.Accept()
		if err != nil {
			incoming <- nil
			return
		}
		log.Printf("[music] Device connected!")

		var buf = make([]byte, 1)
		go func() {
			_, err = conn.Read(buf)
			log.Printf("[music] Connection closed")
		}()

		music := NewMusic(conn)
		incoming <- music
	}()

	time.Sleep(time.Second * 1)

	log.Printf("[music] Initializating Music Mode...")
	err = c.commander.executeCommand(
		partialCommand{"set_music", params{1, hostIP, port}},
	)
	if err != nil {
		return err, nil
	}
	log.Printf("[music] Music Mode Initialized!")

	music := <-incoming
	if music == nil {
		return errors.New("[music] Connection failed"), nil
	}

	return nil, music
}

// chooseEffect returns effect string Accordingly to given duration value.
// Chooses between "sudden" and "smooth".
func chooseEffect(duration int) (string, error) {
	// if duration != 0 && duration < 30 {
	//	return "", errors.New("Ooops")
	// }

	if duration == 0 {
		return "sudden", nil
	}
	return "smooth", nil
}
