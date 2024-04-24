package cloudcontrol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jesper-nord/go-pcc/types"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is a Panasonic Comfort Cloud client.
type Client struct {
	Utoken     string
	DeviceGUID string
	Server     string
}

// SetDevice sets the device GUID on the client.
func (c *Client) SetDevice(deviceGUID string) {
	c.DeviceGUID = deviceGUID
}

// NewClient creates a new Panasonic Comfort Cloud client.
func NewClient() Client {
	return NewClientWithUrl(types.BaseServerUrl)
}

// NewClientWithUrl creates a new client with given base URL.
func NewClientWithUrl(url string) Client {
	client := Client{}
	client.Server = url

	log.Debugf("Created new client for %s", client.Server)

	return client
}

// ValidateSession checks if the session token is still valid.
func (c *Client) ValidateSession(token string) ([]byte, error) {
	c.Utoken = token
	body, err := c.doGetRequest(types.UrlPathValidate)
	if err != nil {
		return body, fmt.Errorf("error: %v %s", err, body)
	}

	return body, nil
}

// CreateSession initialises a client session to Panasonic Comfort Cloud.
func (c *Client) CreateSession(username string, password string) ([]byte, error) {
	postBody, _ := json.Marshal(map[string]any{
		"language": 0,
		"loginId":  username,
		"password": password,
	})

	body, err := c.doPostRequest(types.UrlPathLogin, postBody)
	if err != nil {
		return nil, fmt.Errorf("error: %v %s", err, body)
	}

	session := types.Session{}
	err = json.Unmarshal(body, &session)
	if err != nil {
		log.Fatalf("unmarshal error %v: %s", err, body)
	}

	c.Utoken = session.Utoken

	return body, nil
}

// GetGroups gets all Panasonic Comfort Cloud groups associated to this account.
func (c *Client) GetGroups() (types.Groups, error) {
	body, err := c.doGetRequest(types.UrlPathGroups)
	if err != nil {
		return types.Groups{}, fmt.Errorf("error: %v %s", err, body)
	}
	groups := types.Groups{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		log.Fatalf("unmarshal error %v: %s", err, body)
	}

	return groups, nil
}

// ListDevices lists all available devices.
func (c *Client) ListDevices() ([]string, error) {
	var available []string
	groups, err := c.GetGroups()
	if err != nil {
		return nil, err
	}
	for _, group := range groups.Groups {
		for _, device := range group.Devices {
			available = append(available, device.DeviceGUID)
		}
	}
	return available, nil
}

// GetDeviceStatus gets all details for a specific device.
func (c *Client) GetDeviceStatus() (types.Device, error) {
	body, err := c.doGetRequest(types.UrlPathDeviceStatus + url.QueryEscape(c.DeviceGUID))
	if err != nil {
		return types.Device{}, fmt.Errorf("error: %v %s", err, body)
	}

	device := types.Device{}
	err = json.Unmarshal(body, &device)
	if err != nil {
		log.Fatalf("unmarshal error %v: %s", err, body)
	}

	return device, nil
}

// GetDeviceHistory will fetch historical device data from Panasonic.
func (c *Client) GetDeviceHistory(timeFrame int64) (types.History, error) {
	postBody, _ := json.Marshal(map[string]string{
		"dataMode":   fmt.Sprint(timeFrame),
		"date":       time.Now().Format("20060102"),
		"deviceGuid": c.DeviceGUID,
		"osTimezone": "+01:00",
	})

	body, err := c.doPostRequest(types.UrlPathHistory, postBody)
	if err != nil {
		return types.History{}, fmt.Errorf("error: %v %s", err, body)
	}

	history := types.History{}
	err = json.Unmarshal(body, &history)
	if err != nil {
		log.Fatalf("unmarshal error %v: %s", err, body)
	}

	return history, nil
}

// SetTemperature will set the temperature for a device.
func (c *Client) SetTemperature(temperature float64) ([]byte, error) {
	command := types.Command{
		DeviceGUID: c.DeviceGUID,
		Parameters: types.DeviceControlParameters{
			TemperatureSet: &temperature,
		},
	}

	return c.control(command)
}

// SetFanSpeed will set the fan speed for a device.
func (c *Client) SetFanSpeed(fanSpeed int64) ([]byte, error) {
	command := types.Command{
		DeviceGUID: c.DeviceGUID,
		Parameters: types.DeviceControlParameters{
			FanSpeed: &fanSpeed,
		},
	}

	return c.control(command)
}

// TurnOn will switch the device on.
func (c *Client) TurnOn() ([]byte, error) {
	var on int64 = 1
	command := types.Command{
		DeviceGUID: c.DeviceGUID,
		Parameters: types.DeviceControlParameters{
			Operate: &on,
		},
	}

	return c.control(command)
}

// TurnOff will switch the device off.
func (c *Client) TurnOff() ([]byte, error) {
	var off int64 = 0
	command := types.Command{
		DeviceGUID: c.DeviceGUID,
		Parameters: types.DeviceControlParameters{
			Operate: &off,
		},
	}

	return c.control(command)
}

// SetMode will set the device to the requested AC mode.
func (c *Client) SetMode(mode int64) ([]byte, error) {
	command := types.Command{
		DeviceGUID: c.DeviceGUID,
		Parameters: types.DeviceControlParameters{},
	}

	command.Parameters.OperationMode = &mode

	return c.control(command)
}

// SetEcoMode will set the device to the requested eco mode.
func (c *Client) SetEcoMode(mode int64) ([]byte, error) {
	command := types.Command{
		DeviceGUID: c.DeviceGUID,
		Parameters: types.DeviceControlParameters{},
	}

	command.Parameters.EcoMode = &mode

	return c.control(command)
}

// control sends commands to the Panasonic cloud to control a device.
func (c *Client) control(command types.Command) ([]byte, error) {
	postBody, _ := json.Marshal(command)

	log.Debugf("Command: %s", postBody)

	body, err := c.doPostRequest(types.UrlPathControl, postBody)
	if err != nil {
		return nil, fmt.Errorf("error: %v %s", err, body)
	}
	if string(body) != types.SuccessResponse {
		return body, fmt.Errorf("error body: %v %s", err, body)
	}

	return body, nil
}

func (c *Client) doPostRequest(url string, postbody []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.Server+url, bytes.NewBuffer(postbody))
	c.setHeaders(req)

	log.Debugf("POST request URL: %#v\n", req.URL)
	log.Debugf("POST request body: %#v\n", string(postbody))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	log.Debugf("POST response body: %s", string(body))

	if resp.StatusCode > 200 {
		return body, fmt.Errorf("HTTP Error: %s", resp.Status)
	}

	return body, nil
}

func (c *Client) doGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.Server+url, nil)
	c.setHeaders(req)

	log.Debugf("GET request URL: %#v", req.URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	log.Debugf("GET response body: %s", string(body))

	if resp.StatusCode > 200 {
		return body, fmt.Errorf("HTTP Error: %s", resp.Status)
	}

	return body, nil
}

func (c *Client) setHeaders(req *http.Request) {
	if c.Utoken != "" {
		req.Header.Set("X-User-Authorization", c.Utoken)
	}
	req.Header.Set("X-APP-TYPE", "1")
	req.Header.Set("X-APP-VERSION", "1.20.0")
	req.Header.Set("X-APP-TIMESTAMP", "1")
	req.Header.Set("X-APP-NAME", "Comfort Cloud")
	req.Header.Set("X-CFC-API-KEY", "Comfort Cloud")
	req.Header.Set("User-Agent", "G-RAC")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "Keep-Alive")

	log.Debugf("HTTP headers set to: %#v", req.Header)
}
