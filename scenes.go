package yeelight

type Scene interface {
	toParams() params
}

type BaseScene struct {
	name string
}

type ColorScene struct {
	BaseScene
	rgb, brightness int
}

func (s ColorScene) toParams() params {
	return params{s.name, s.rgb, s.brightness}
}

type HSVScene struct {
	BaseScene
	hue, saturation, brightness int
}

func (s HSVScene) toParams() params {
	return params{s.name, s.hue, s.saturation, s.brightness}
}

type TemperatureScene struct {
	BaseScene
	temperature, brightness int
}

func (s TemperatureScene) toParams() params {
	return params{s.name, s.temperature, s.brightness}
}

type ColorFlowScene struct {
	BaseScene
	count, action  int
	flowExpression FlowExpression
}

func (s ColorFlowScene) toParams() params {
	return params{s.name, s.count, s.action, s.flowExpression}
}

// automatic shutdown after specified amount of minutes
type AutoDelayOffScene struct {
	BaseScene
	brightness, minutes int
}

func (s AutoDelayOffScene) toParams() params {
	return params{s.name, s.brightness, s.minutes}
}
