package dvrip

import (
	"encoding/json"
	"errors"
	"fmt"
)

type SystemInfoResponse struct {
	*Response
	Info systemInfo `json:"SystemInfo"`
}

type systemInfo struct {
	AlarmInChannel  int    `json:"AlarmInChannel"`
	AlarmOutChannel int    `json:"AlarmOutChannel"`
	AudioInChannel  int    `json:"AudioInChannel"`
	BuildTime       string `json:"BuildTime"`
	CombineSwitch   int    `json:"CombineSwitch"`
	DeviceRunTime   string `json:"DeviceRunTime"`
	DeviceType      int    `json:"DeviceType"`
	DigChannel      int    `json:"DigChannel"`
	EncryptVersion  string `json:"EncryptVersion"`
	ExtraChannel    int    `json:"ExtraChannel"`
	HardWare        string `json:"HardWare"`
	HardWareVersion string `json:"HardWareVersion"`
	SerialNo        string `json:"SerialNo"`
	SoftWareVersion string `json:"SoftWareVersion"`
	TalkInChannel   int    `json:"TalkInChannel"`
	TalkOutChannel  int    `json:"TalkOutChannel"`
	UpdataTime      string `json:"UpdataTime"`
	UpdataType      string `json:"UpdataType"`
	VideoInChannel  int    `json:"VideoInChannel"`
	VideoOutChannel int    `json:"VideoOutChannel"`
}

func (c *Client) GetSystemInfo() (*SystemInfoResponse, error) {
	var result = new(SystemInfoResponse)

	_, resp, err := c.Instruct(CodeSystemInfo, "", nil)
	if err != nil {
		return nil, errors.New("failed to send instruction to device")
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	if result.Ret != statusOK {
		return nil, fmt.Errorf("unexpected status code; %s", mappedStatusCodes[result.Ret])
	}

	return result, nil
}
