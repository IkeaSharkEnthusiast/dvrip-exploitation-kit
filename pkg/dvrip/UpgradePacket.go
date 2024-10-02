package dvrip

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const (
	upgradePacketSize = 0x8000
)

type headerUpgrade struct {
	Head    uint8 // B
	Version uint8 // B

	_ [2]byte // 2x

	Session        uint32 // I
	SequenceNumber uint32 // I

	_ byte // x

	Byte3 byte   // B
	MsgID uint16 // H
	Len   uint32 // I
}

func (c *Client) Upgrade(fileData []byte) error {
	var data = make(map[string]interface{})
	var filePacket = make([]byte, upgradePacketSize)

	_, bytesR, err := c.Instruct(CodeStartSystemUpgrade, "", map[string]interface{}{
		"Action": "Start",
		"Type":   "System",
	})

	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytesR, &data); err != nil {
		return err
	}

	if data["Ret"].(float64) != float64(statusOK) {
		return fmt.Errorf("unexpected status code: %s", mappedStatusCodes[int(data["Ret"].(float64))])
	}

	var blockNum uint32 = 0
	var scanner = bytes.NewReader(fileData)

	for {
		n, err := scanner.Read(filePacket)
		if err != nil || err == io.EOF {
			if err == io.EOF {
				break
			}

			return err
		}

		_, bytesR, err = c.InstructRaw(CodeSendFile, blockNum, filePacket[:n], uint32(n))
		if err != nil {
			return fmt.Errorf("failed to send data: %v", err)
		}

		blockNum++

		if err := json.Unmarshal(bytesR, &data); err != nil {
			return fmt.Errorf("failed to parse data response: %v", err)
		}

		if data["Ret"].(float64) != float64(statusOK) {
			return fmt.Errorf("unexpected status code: %s", mappedStatusCodes[int(data["Ret"].(float64))])
		}
	}

	err = c.writeUpgradeHeader(blockNum)
	if err != nil {
		return err
	}

	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// read the response of the actual start upgrade header
	_, bytesR, err = c.Read(true)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytesR, &data); err != nil {
		return err
	}

	if !(data["Ret"].(float64) == float64(statusOK) || data["Ret"].(float64) == float64(statusUpgradeSuccessful)) {
		return fmt.Errorf("unexpected status code: %s", mappedStatusCodes[int(data["Ret"].(float64))])
	}

	// now read the actual response
	_, bytesR, err = c.Read(true)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytesR, &data); err != nil {
		return err
	}

	if !(data["Ret"].(float64) == float64(statusOK) || data["Ret"].(float64) == float64(statusUpgradeSuccessful)) {
		return fmt.Errorf("unexpected status code: %s", mappedStatusCodes[int(data["Ret"].(float64))])
	}

	return nil
}

func (c *Client) writeUpgradeHeader(sequence uint32) error {
	var buf bytes.Buffer

	// send the header
	if err := binary.Write(&buf, binary.LittleEndian, headerUpgrade{
		Head:           255,
		Version:        0,
		Session:        c.Session,
		SequenceNumber: sequence,
		Byte3:          1,
		MsgID:          0x05F2,
		Len:            0,
	}); err != nil {
		return err
	}

	// send all the data to the device
	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
