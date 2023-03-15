package types

// HistoryDataMode maps out the time intervals to fetch history data
var HistoryDataMode = map[string]int64{
	"day":   0,
	"week":  1,
	"month": 2,
	"year":  3,
}

// Modes define the different AC modes the device can be in
var Modes = map[string]int64{
	"auto": 0,
	"dry":  1,
	"cool": 2,
	"heat": 3,
	"fan":  4,
}

// ModesReverse define the different AC modes the device can be in
var ModesReverse = map[int64]string{
	0: "auto",
	1: "dry",
	2: "cool",
	3: "heat",
	4: "fan",
}

var EcoMode = map[string]int64{
	"auto":     0,
	"powerful": 1,
	"quiet":    2,
}

var EcoModeReverse = map[int64]string{
	0: "auto",
	1: "powerful",
	2: "quiet",
}

var FanSpeedReverse = map[int64]string{
	0: "auto",
	1: "1",
	2: "2",
	3: "3",
	4: "4",
	5: "5",
}

// Operate defines if the AC is on or off
var Operate = map[int64]string{
	0: "off",
	1: "on",
}

// Session is a login session structure
type Session struct {
	Utoken   string `json:"uToken"`
	Result   int64  `json:"result"`
	Language int64  `json:"language"`
}

// Groups is a set of grouped devices
type Groups struct {
	GroupCount int64   `json:"groupCount"`
	Groups     []Group `json:"groupList"`
}

// Group defines a control group with devices
type Group struct {
	GroupID   int64    `json:"groupId"`
	GroupName string   `json:"groupName"`
	Devices   []Device `json:"deviceList"`
}

// DeviceControlParameters are the device control parameters
// used in Marshalling control commands
// We need to duplicate this with pointers to make sure the
// 'omitempty' parameter will not cancel out eg operate = 0
// when sending control commands to the unit
type DeviceControlParameters struct {
	ActualNanoe             *int64   `json:"actualNanoe,omitempty"`
	AirDirection            *int64   `json:"airDirection,omitempty"`
	AirQuality              *int64   `json:"airQuality,omitempty"`
	AirSwingLR              *int64   `json:"airSwingLR,omitempty"`
	AirSwingUD              *int64   `json:"airSwingUD,omitempty"`
	Defrosting              *int64   `json:"defrosting,omitempty"`
	DevGUID                 *string  `json:"devGuid,omitempty"`
	DevRacCommunicateStatus *int64   `json:"devRacCommunicateStatus,omitempty"`
	EcoFunctionData         *int64   `json:"ecoFunctionData,omitempty"`
	EcoMode                 *int64   `json:"ecoMode,omitempty"`
	EcoNavi                 *int64   `json:"ecoNavi,omitempty"`
	Permission              *int64   `json:"permission,omitempty"`
	ErrorCode               *int64   `json:"errorCode,omitempty"`
	ErrorCodeStr            *string  `json:"errorCodeStr,omitempty"`
	ErrorStatus             *int64   `json:"errorStatus,omitempty"`
	ErrorStatusFlg          *bool    `json:"errorStatusFlg,omitempty"`
	FanAutoMode             *int64   `json:"fanAutoMode,omitempty"`
	FanSpeed                *int64   `json:"fanSpeed,omitempty"`
	HTTPErrorCode           *int64   `json:"httpErrorCode,omitempty"`
	Iauto                   *int64   `json:"iAuto,omitempty"`
	InsideTemperature       *float64 `json:"insideTemperature,omitempty"`
	Nanoe                   *int64   `json:"nanoe,omitempty"`
	Online                  *bool    `json:"online,omitempty"`
	Operate                 *int64   `json:"operate,omitempty"`       // Turn on/off
	OperationMode           *int64   `json:"operationMode,omitempty"` // Set Mode (heat, dry, etc)
	OutsideTemperature      *float64 `json:"outTemperature,omitempty"`
	PowerfulMode            *bool    `json:"powerfulMode,omitempty"`
	TemperatureSet          *float64 `json:"temperatureSet,omitempty"` // Set Temperature
	UpdateTime              *int64   `json:"updateTime,omitempty"`
}

// DeviceParameters are the current device parameters
// Used when UnMarshalling current device status
type DeviceParameters struct {
	ActualNanoe             int64   `json:"actualNanoe"`
	AirDirection            int64   `json:"airDirection"`
	AirQuality              int64   `json:"airQuality"`
	AirSwingLR              int64   `json:"airSwingLR"`
	AirSwingUD              int64   `json:"airSwingUD"`
	Defrosting              int64   `json:"defrosting"`
	DevGUID                 string  `json:"devGuid"`
	DevRacCommunicateStatus int64   `json:"devRacCommunicateStatus"`
	EcoFunctionData         int64   `json:"ecoFunctionData"`
	EcoMode                 int64   `json:"ecoMode"`
	EcoNavi                 int64   `json:"ecoNavi"`
	Permission              int64   `json:"permission"`
	ErrorCode               int64   `json:"errorCode"`
	ErrorCodeStr            string  `json:"errorCodeStr"`
	ErrorStatus             int64   `json:"errorStatus"`
	ErrorStatusFlg          bool    `json:"errorStatusFlg"`
	FanAutoMode             int64   `json:"fanAutoMode"`
	FanSpeed                int64   `json:"fanSpeed"`
	HTTPErrorCode           int64   `json:"httpErrorCode"`
	Iauto                   int64   `json:"iAuto"`
	InsideTemperature       float64 `json:"insideTemperature"`
	Nanoe                   int64   `json:"nanoe"`
	Online                  bool    `json:"online"`
	Operate                 int64   `json:"operate"`       // on/off
	OperationMode           int64   `json:"operationMode"` // Mode (heat, dry, etc)
	OutsideTemperature      float64 `json:"outTemperature"`
	PowerfulMode            bool    `json:"powerfulMode"`
	TemperatureSet          float64 `json:"temperatureSet"` // Temperature
	UpdateTime              int64   `json:"updateTime"`
}

// Device is Panasonic device
type Device struct {
	AirSwingLR         bool             `json:"airSwingLR"`
	AutoMode           bool             `json:"autoMode"`
	AutoTempMax        int64            `json:"autoTempMax"`
	AutoTempMin        int64            `json:"autoTempMin"`
	CoolMode           bool             `json:"coolMode"`
	CoolTempMax        int64            `json:"coolTeampMax"`
	CoolTempMin        int64            `json:"coolTempMin"`
	DeviceGUID         string           `json:"deviceGuid"`
	DeviceHashGUID     string           `json:"deviceHashGuid"`
	DeviceModuleNumber string           `json:"deviceModuleNumber"`
	DeviceName         string           `json:"deviceName"`
	DeviceType         string           `json:"deviceType"`
	DryMode            bool             `json:"dryMode"`
	DryTempMax         int64            `json:"dryTempMax"`
	DryTempMin         int64            `json:"dryTempMin"`
	EcoFunction        int64            `json:"ecoFunction"`
	EcoNavi            bool             `json:"ecoNavi"`
	FanDirectionMode   int64            `json:"fanDirectionMode"`
	FanMode            bool             `json:"fanMode"`
	FanSpeedMode       int64            `json:"fanSpeedMode"`
	HeatMode           bool             `json:"heatMode"`
	HeatTempMax        int64            `json:"heatTempMax"`
	HeatTempMin        int64            `json:"heatTeampMin"`
	IautoX             bool             `json:"iAutoX"`
	ModeAvlAutoMode    bool             `json:"modeAvlList.autoMode"`
	ModeAvlFanMode     bool             `json:"modeAvlList.fanMode"`
	Nanoe              bool             `json:"nanoe"`
	QuietMode          bool             `json:"quietMode"`
	SummerHouse        int64            `json:"summerHouse"`
	TemperatureUnit    int64            `json:"temperatureUnit"`
	TimeStamp          int64            `json:"timestamp"`
	Parameters         DeviceParameters `json:"parameters"`
}

// History is a list of HistoryEntry points with measurements
type History struct {
	EnergyConsumption  float64        `json:"energyConsumption"`
	EstimatedCost      float64        `json:"estimatedCost"`
	DeviceRegisterTime string         `json:"deviceRegisterTime"`
	CurrencyUnit       string         `json:"currencyUnit"`
	HistoryEntries     []HistoryEntry `json:"historyDataList"`
}

// HistoryEntry is detailed data for a given day,week,month,year
type HistoryEntry struct {
	DataNumber         int64   `json:"dataNumber"`
	Consumption        float64 `json:"consumption"`
	Cost               float64 `json:"cost"`
	AverageSettingTemp float64 `json:"averageSettingTemp"`
	AverageInsideTemp  float64 `json:"averageInsideTemp"`
	AverageOutsideTemp float64 `json:"averageOutsideTemp"`
}

// Command is basic command control structure
type Command struct {
	DeviceGUID string                  `json:"deviceGuid"`
	Parameters DeviceControlParameters `json:"parameters"`
}
