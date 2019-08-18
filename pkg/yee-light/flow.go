package yeelight

import (
	"errors"
	"fmt"
)

type FlowState struct {
	Duration int    // standard transition duration in milliseconds (minimum 50)
	Mode     CfMode // 1: color, 2: temperature, 7: sleep

	// rgb value (0-0xffffff) for color mode
	// temperature value (1700-6500) for temperature mode
	// milliseconds (?-?) for sleep mode
	Value int

	Brightness int // brightness value (1-100), -1 when don't want to change brightness
}

// NewFlowState creates transition step and panics on incorrect input variables
// examples:
//   yl.NewFlowState(50, yl.CF_MODE_COLOR, 0xff0000, 100)
//   yl.NewFlowState(200, yl.CF_MODE_SLEEP, 0, 0)
//   yl.NewFlowState(50, yl.CF_MODE_COLOR, 0x0000ff, 100)
//   yl.NewFlowState(200, yl.CF_MODE_SLEEP, 0, 0)
// TODO: prepare unfriendly interface for creating transition objects separately, according to supported modes
//      for instance, current sleep command:
//          NewFlowState(200, yl.CF_MODE_SLEEP, 0, 0)
//      could be changed to:
//          NewSleepState(200)
func NewFlowState(duration int, mode CfMode, value, brightness int) FlowState {
	if duration < 50 {
		panic("duration required to be >= 50")
	}

	valudateBrightness := func(brightness int) {
		if brightness < 1 || brightness > 100 {
			// documentation says -1 is possible for skipping brightness manipulation,
			// however, It just doesn't work, at least for my bulb (general error is returned with code 5000)
			// according to this founding, I'll stick with standard 1-100 range validation
			panic("brightness in 1-100 range or -1 (do not change brightness)")
		}
	}

	switch mode {
	case CF_MODE_COLOR:
		if value < 0 || value > 0xffffff {
			panic("value for color mode should be in 0-0xffffff range")
		}
		valudateBrightness(brightness)
	case CF_MODE_TEMP:
		if value < 1700 || value > 6500 {
			panic("value for temperature mode should be in 1700-6500 range")
		}
		valudateBrightness(brightness)
	case CF_MODE_SLEEP:
		// value and brightness are ignored in sleep mode
		value = 0
		brightness = 0
	default:
		panic("mode required to be 1 (color), 2 (temp), or 7 (sleep)")
	}

	return FlowState{duration, mode, value, brightness}
}

type FlowExpression struct {
	states []FlowState
}

func (e *FlowExpression) encode() string {
	var encodedExpression string
	var lastElement int = len(e.states) - 1

	for i, state := range e.states {
		stateEncoded := fmt.Sprintf(
			"%d,%d,%d,%d", state.Duration, state.Mode, state.Value, state.Brightness,
		)
		encodedExpression += stateEncoded

		if i != lastElement {
			encodedExpression += ","
		}
	}
	return encodedExpression
}

func NewFlowExpression(states ...FlowState) (error, FlowExpression) {
	if len(states) == 0 {
		return errors.New("flowExpression should have at least one FlowState, please pass one"), FlowExpression{}
	}

	var flowStates []FlowState

	for _, state := range states {
		flowStates = append(flowStates, state)
	}

	return nil, FlowExpression{states: states}
}
