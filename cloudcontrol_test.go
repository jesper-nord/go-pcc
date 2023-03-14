package cloudcontrol_test

import (
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	cloudcontrol "github.com/hacktobeer/go-panasonic"
	pt "github.com/hacktobeer/go-panasonic/types"
)

var (
	client      cloudcontrol.Client
	sessionBody = `{"uToken":"token12345","language":0,"result":0}`
	groupsBody  = `{"iaqStatus":{"statusCode":200},"groupCount":1,"groupList":[{"groupId":112867,"groupName":"My House","deviceList":[{"deviceGuid":"CZ-CAPWFC1+B8B7F1B3E326","deviceType":"4","deviceName":"Alaior-home","permission":3,"deviceModuleNumber":"S-125PU2E5B","deviceHashGuid":"f609023332bbeee157a5b868fe80b9fb14a1d883938c1836003796332150db16","summerHouse":0,"iAutoX":false,"nanoe":true,"autoMode":true,"heatMode":true,"fanMode":false,"dryMode":true,"coolMode":true,"ecoNavi":false,"powerfulMode":true,"quietMode":true,"airSwingLR":true,"ecoFunction":0,"temperatureUnit":0,"modeAvlList":{"autoMode":1,"fanMode":1},"autoTempMax":27,"autoTempMin":17,"dryTempMax":30,"dryTempMin":18,"coolTempMax":30,"coolTempMin":18,"heatTempMax":30,"heatTempMin":16,"fanSpeedMode":5,"fanDirectionMode":5,"parameters":{"operate":1,"operationMode":0,"temperatureSet":19.5,"fanSpeed":0,"fanAutoMode":1,"airSwingLR":2,"airSwingUD":3,"ecoMode":0,"ecoNavi":0,"nanoe":1,"iAuto":0,"actualNanoe":1,"airDirection":3,"ecoFunctionData":0}}]}]}`
	historyBody = `{"energyConsumption":2.9,"estimatedCost":0.0,"deviceRegisterTime":"20201216","currencyUnit":"€","historyDataList":[{"dataNumber":0,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":23.0,"averageOutsideTemp":14.0},{"dataNumber":1,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":23.0,"averageOutsideTemp":13.75},{"dataNumber":2,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":23.0,"averageOutsideTemp":13.0},{"dataNumber":3,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":23.0,"averageOutsideTemp":13.0},{"dataNumber":4,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":22.75,"averageOutsideTemp":13.0},{"dataNumber":5,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":22.0,"averageOutsideTemp":13.0},{"dataNumber":6,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.75,"averageInsideTemp":20.75,"averageOutsideTemp":12.75},{"dataNumber":7,"consumption":0.5,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":18.75,"averageOutsideTemp":11.25},{"dataNumber":8,"consumption":0.4,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":18.0,"averageOutsideTemp":12.25},{"dataNumber":9,"consumption":0.3,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":20.75,"averageOutsideTemp":13.75},{"dataNumber":10,"consumption":0.2,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":23.0,"averageOutsideTemp":14.0},{"dataNumber":11,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":23.0,"averageOutsideTemp":14.5},{"dataNumber":12,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":23.0,"averageOutsideTemp":15.0},{"dataNumber":13,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":21.0,"averageOutsideTemp":15.25},{"dataNumber":14,"consumption":0.4,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":18.0,"averageOutsideTemp":15.5},{"dataNumber":15,"consumption":0.2,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":18.5,"averageOutsideTemp":16.0},{"dataNumber":16,"consumption":0.2,"cost":0.0,"averageSettingTemp":19.0,"averageInsideTemp":19.0,"averageOutsideTemp":15.0},{"dataNumber":17,"consumption":0.2,"cost":0.0,"averageSettingTemp":19.125,"averageInsideTemp":19.0,"averageOutsideTemp":14.25},{"dataNumber":18,"consumption":0.2,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":18.75,"averageOutsideTemp":13.5},{"dataNumber":19,"consumption":0.3,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":18.75,"averageOutsideTemp":12.0},{"dataNumber":20,"consumption":0.0,"cost":0.0,"averageSettingTemp":19.5,"averageInsideTemp":19.0,"averageOutsideTemp":11.0},{"dataNumber":21,"consumption":-255,"cost":-255,"averageSettingTemp":-255,"averageInsideTemp":-255,"averageOutsideTemp":-255},{"dataNumber":22,"consumption":-255,"cost":-255,"averageSettingTemp":-255,"averageInsideTemp":-255,"averageOutsideTemp":-255},{"dataNumber":23,"consumption":-255,"cost":-255,"averageSettingTemp":-255,"averageInsideTemp":-255,"averageOutsideTemp":-255}],"temperatureUnit":0}`
	controlBody = `{"result":0}`
)

func TestMain(m *testing.M) {
	server := serverMock()
	defer server.Close()

	client = cloudcontrol.NewClientWithUrl(server.URL)

	os.Exit(m.Run())
}

func TestNewClient(t *testing.T) {
	client := cloudcontrol.NewClient()
	got := client.Server

	assert.Equal(t, got, pt.URLServer)
}

func TestSetDevice(t *testing.T) {
	device := "device12345"

	var client cloudcontrol.Client
	client.SetDevice(device)

	want := device
	got := client.DeviceGUID

	assert.Equal(t, got, want)
}

func TestTurnOn(t *testing.T) {
	client.CreateSession("", "")
	body, err := client.TurnOn()
	if err != nil {
		t.Errorf("TestTurnOn() returned an error: %v", err)
	}

	want := pt.SuccessResponse
	got := string(body)

	assert.Equal(t, got, want)
}

func TestGetGroups(t *testing.T) {
	client.CreateSession("", "")
	groups, _ := client.GetGroups()

	want := "My House"
	got := groups.Groups[0].GroupName

	assert.Equal(t, got, want)
	assert.Equal(t, 1, len(groups.Groups[0].Devices))
}

func TestGetDeviceHistory(t *testing.T) {
	client.CreateSession("", "")
	history, err := client.GetDeviceHistory(pt.HistoryDataMode["day"])
	if err != nil {
		t.Error(err)
	}

	got := len(history.HistoryEntries)
	want := 24

	assert.Equal(t, got, want)
}

func TestCreateSession(t *testing.T) {
	username := "test@test.com"
	password := "secret1234"

	client.CreateSession(username, password)

	got := client.Utoken
	want := "token12345"

	assert.Equal(t, got, want)
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc(pt.URLLogin, sessionMock)
	handler.HandleFunc(pt.URLGroups, groupsMock)
	handler.HandleFunc(pt.URLControl, controlMock)
	handler.HandleFunc(pt.URLHistory, historyMock)

	srv := httptest.NewServer(handler)

	return srv
}

func historyMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(historyBody))
}

func sessionMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(sessionBody))
}

func groupsMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(groupsBody))
}

func controlMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(controlBody))
}
