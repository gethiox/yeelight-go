package yeelight

import (
	"errors"
	"time"
)

// Common commands defines shared commands for standard mode and background mode
type commonCommands struct {
	commander commander
	prefix    string // used as a prefix in command names, empty by default. can support background commands by "bg_".
}

// Temperature sets device temperature, range 1700-6500
func (c *commonCommands) Temperature(temp int, duration time.Duration) error {
	if temp < 1700 || temp > 6500 {
		return errors.New("temperature expected range: 1700-6500")
	}

	effect, err := chooseEffect(duration)
	if err != nil {
		return err
	}

	return c.commander.executeCommand(
		partialCommand{c.prefix + "set_ct_abx", params{temp, effect, timeToMs(duration)}},
	)
}

// RGB sets device color in RGB form. range 0x000000-0xFFFFFF
// It's easy to prepare RGB as hexadecimal value, color order: 0xRRGGBB
// example: sets color to red, and then blue
//   bulb.RGB(0xff0000, 0)
//   bulb.RGB(0x0000ff, 0)
func (c *commonCommands) RGB(rgb int, duration time.Duration) error {
	if rgb < 0 || rgb > 0xffffff {
		return errors.New("rgb expected range: 0-0xFFFFFF")
	}

	effect, err := chooseEffect(duration)
	if err != nil {
		return err
	}

	return c.commander.executeCommand(
		partialCommand{c.prefix + "set_rgb", params{rgb, effect, timeToMs(duration)}},
	)
}

// HSV sets device color in HSV form. hue range: 0-359, saturation range: 0-100
func (c *commonCommands) HSV(hue, saturation int, duration time.Duration) error {
	if hue < 0 || hue > 359 {
		return errors.New("hue expected range: 0-359")
	}

	if saturation < 0 || saturation > 100 {
		return errors.New("saturation expected range: 0-100")
	}

	effect, err := chooseEffect(duration)
	if err != nil {
		return err
	}

	return c.commander.executeCommand(
		partialCommand{c.prefix + "set_hsv", params{hue, saturation, effect, timeToMs(duration)}},
	)
}

// Brightness sets device brightness in range 1-100
func (c *commonCommands) Brightness(brightness int, duration time.Duration) error {
	if brightness < 1 || brightness > 100 {
		return errors.New("brightness expected range: 1-100")
	}

	effect, err := chooseEffect(duration)
	if err != nil {
		return err
	}

	return c.commander.executeCommand(
		partialCommand{c.prefix + "set_bright", params{brightness, effect, timeToMs(duration)}},
	)
}

// StartColorFlow sets device in color flow mode, FlowExpression determines wanted animation..
// It can be changing brightness, color or temperature.
func (c *commonCommands) StartColorFlow(count int, action CfAction, flowExpression FlowExpression) error {
	return c.commander.executeCommand(
		partialCommand{c.prefix + "start_cf", params{count, action, flowExpression.encode()}},
	)
}

func (c *commonCommands) StopColorFlow() error {
	return c.commander.executeCommand(
		partialCommand{c.prefix + "stop_cf", params{}},
	)
}

// SetScene can change state to given Scene, even if current device state is "off"
// Not implemented! TODO: IMPLEMENT
func (c *commonCommands) SetScene(scene Scene) error {
	return errors.New("not implemented")
}

// Sets current state as default
func (c *commonCommands) SetDefault() error {
	return c.commander.executeCommand(
		partialCommand{c.prefix + "set_default", params{}},
	)
}

func (c *commonCommands) PowerOn(duration time.Duration) error {
	effect, err := chooseEffect(duration)
	if err != nil {
		return err
	}

	return c.commander.executeCommand(
		partialCommand{c.prefix + "set_power", params{"on", effect, timeToMs(duration)}},
	)
}

// PowerOnWithMode behaves similarly to ordinal PowerOn except it can sets device directly to given Mode
func (c *commonCommands) PowerOnWithMode(duration time.Duration, mode Mode) error {
	effect, err := chooseEffect(duration)
	if err != nil {
		return err
	}

	return c.commander.executeCommand(
		partialCommand{c.prefix + "set_power", params{"on", effect, timeToMs(duration), int(mode)}},
	)
}

func (c *commonCommands) PowerOff(duration time.Duration) error {
	effect, err := chooseEffect(duration)
	if err != nil {
		return err
	}

	return c.commander.executeCommand(
		partialCommand{c.prefix + "set_power", params{"off", effect, timeToMs(duration)}},
	)
}

// Toggle is a Built-in method which toggles device state.
// Only limitation is that fade effect can't be modified here, use PowerOn and PowerOff instead
func (c *commonCommands) Toggle() error {
	return c.commander.executeCommand(
		partialCommand{c.prefix + "toggle", params{}},
	)
}

type backgroundLightCommands struct {
	commander commander
	commonCommands
	prefix string // required to be set to "bg_" !!!
}

// DevToggle is toggling the main light and background light at the same time
func (c *backgroundLightCommands) DevToggle() error {
	return c.commander.executeCommand(
		partialCommand{"dev_toggle", params{}},
	)
}
