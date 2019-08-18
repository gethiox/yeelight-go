package yeelight

const (
	LF = byte('\n')
	CR = byte('\r')

	// Properties
	PROP_POWER          Property = "power"          // on: smart LED is turned on  /  off: smart LED is turned off
	PROP_BRIGHT         Property = "bright"         // Brightness percentage. Range 1~100
	PROP_CT             Property = "ct"             // Color temperature. Range 1700~6500(k)
	PROP_RGB            Property = "rgb"            // Color. Range 1~16777215
	PROP_HUE            Property = "hue"            // Hue. Range 0~359
	PROP_SAT            Property = "sat"            // Saturation. Range 0~100
	PROP_COLOR_MODE     Property = "color_mode"     // 1: rgb mode / 2: color temperature mode / 3: hsv mode
	PROP_FLOWING        Property = "flowing"        // 0: no flow is running / 1:color flow is running
	PROP_DELAYOFF       Property = "delayoff"       // The remaining time of a sleep timer. Range 1~60(minutes)
	PROP_FLOW_PARAMS    Property = "flow_params"    // Current flow parameters (only meaningful when 'flowing' is 1)
	PROP_MUSIC_ON       Property = "music_on"       // 1: Music mode is on / 0: Music mode is off
	PROP_NAME           Property = "name"           // The name of the device set by “set_name” Command
	PROP_BG_POWER       Property = "bg_power"       // Background light power status
	PROP_BG_FLOWING     Property = "bg_flowing"     // Background light is flowing
	PROP_BG_FLOW_PARAMS Property = "bg_flow_params" // Current flow parametersof background light
	PROP_BG_CT          Property = "bg_ct"          // Color temperatureof background light
	PROP_BG_LMODE       Property = "bg_lmode"       // 1: rgb mode / 2: color temperature mode / 3: hsv mode
	PROP_BG_BRIGHT      Property = "bg_bright"      // Brightness percentageof background light
	PROP_BG_RGB         Property = "bg_rgb"         // Colorof background light
	PROP_BG_HUE         Property = "bg_hue"         // Hueof background light
	PROP_BG_SAT         Property = "bg_sat"         // Saturationof background light
	PROP_NL_BR          Property = "nl_br"          // Brightness of night mode light
	PROP_ACTIVE_MODE    Property = "active_mode"    // 0: daylight mode / 1: moonlight mode (ceiling light only)

	CRON_TYPE_POWER_OFF CronType = 0 // power off

	ADJUST_ACTION_INCRASE Action     = "incrase"
	ADJUST_ACTION_DECRASE Action     = "decrase"
	ADJUST_ACTION_CIRCLE  Action     = "circle" // incrase
	ADJUST_PROP_BRIGHT    AdjustProp = "bright"
	ADJUST_PROP_CT        AdjustProp = "ct"
	ADJUST_PROP_COLOR     AdjustProp = "color"

	MODE_DEFAUTL Mode = 0 // Normal turn on operation (default value)
	MDOE_CT      Mode = 1 // Turn on and switch to CT mode.
	MODE_RGB     Mode = 2 // Turn on and switch to RGB mode
	MODE_HSV     Mode = 3 // Turn on and switch to HSV mode.
	MODE_CF      Mode = 4 // Turn on and switch to color flow mode
	MODE_NL      Mode = 5 // Turn on and switch to Night light mode. (Ceiling light only).

	CF_MODE_COLOR        CfMode   = 1
	CF_MODE_TEMP         CfMode   = 2
	CF_MODE_SLEEP        CfMode   = 7
	CF_ACTION_RECOVER    CfAction = 0 // smart LED recover to the state before the color flow started.
	CF_ACTION_STAY       CfAction = 1 // smart LED stay at the state when the flow is stopped.
	CF_ACTION_POWEROFF   CfAction = 2 // turn off the smart LED after the flow is stopped.
	CF_COUNT_INF         int      = 0
	CF_BRIGHTNESS_IGNORE int      = -1 // not supported on my device, but documentations says it's supported ¯\_(ツ)_/¯
)

type (
	params     []interface{}
	Mode       int
	CronType   int
	Property   string
	Action     string
	AdjustProp string
	CfAction   int
	CfMode     int
)

type completeCommand struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params params `json:"params"`
}

func newCompleteCommand(partialCommand partialCommand, commandID int) completeCommand {
	return completeCommand{
		ID:     commandID,
		Method: partialCommand.Method,
		Params: partialCommand.Params,
	}
}

type partialCommand struct {
	Method string
	Params params
}

type commander interface {
	executeCommand(partialCommand) error
}
