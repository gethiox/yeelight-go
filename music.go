package yeelight

import (
	"encoding/json"
	"net"
	"time"
)

// Music mode in theory supports all commands, but due to device behaviour
// commands are somehow limited
// for instance you can PowerOff bulb in music mode but as a consequence
// device will also exits music mode immediately ¯\_(ツ)_/¯
// I decided to expose only most useful commands here
// Also in music mode device doesn't respond on commands so errors cannot be returned
type musicSupportedCommands interface {
	Temperature(temp int, duration time.Duration)
	RGB(rgb int, duration time.Duration)
	HSV(hue, saturation int, duration time.Duration)
	Brightness(brightness int, duration time.Duration)
	StartColorFlow(count int, action CfAction, flowExpression FlowExpression)
	StopColorFlow()
}

type Music struct {
	commonCommands

	conn net.Conn
}

func (m *Music) executeCommand(c partialCommand) error {
	// ID doesn't matter in music mode, bulbs doesn't respond an commands
	realCommand := newCompleteCommand(c, 0)
	message, err := json.Marshal(realCommand)
	if err != nil {
		return err
	}
	message = append(message, CR, LF)

	_, err = m.conn.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (m *Music) Stop() error {
	return m.conn.Close()
}

func NewMusic(conn net.Conn) *Music {
	music := &Music{
		commonCommands{},
		conn,
	}
	music.commonCommands.commander = music

	return music
}

func (m *Music) Temperature(temp int, duration time.Duration) {
	_ = m.commonCommands.Temperature(temp, duration)
}

func (m *Music) RGB(rgb int, duration time.Duration) {
	_ = m.commonCommands.RGB(rgb, duration)
}

func (m *Music) HSV(hue, saturation int, duration time.Duration) {
	_ = m.commonCommands.HSV(hue, saturation, duration)
}

func (m *Music) Brightness(brightness int, duration time.Duration) {
	_ = m.commonCommands.Brightness(brightness, duration)
}

func (m *Music) StartColorFlow(count int, action CfAction, flowExpression FlowExpression) {
	_ = m.commonCommands.StartColorFlow(count, action, flowExpression)
}

func (m *Music) StopColorFlow() {
	_ = m.commonCommands.StopColorFlow()
}
