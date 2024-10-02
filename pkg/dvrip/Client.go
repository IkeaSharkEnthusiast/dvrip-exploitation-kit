package dvrip

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func New(target string, timeout time.Duration) (*Client, error) {
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:     conn,
		lastPing: 0,

		Session:  0,
		Sequence: 0,
	}, nil
}

func NewConn(conn net.Conn) (*Client, error) {
	return &Client{
		conn:     conn,
		lastPing: 0,

		Session:  0,
		Sequence: 0,
	}, nil
}

// Write writes something to the device, wow!
func (c *Client) write(msgID uint16, data []byte, dataLen uint32, sequence uint32) error {
	var buf bytes.Buffer

	// send the header
	if err := binary.Write(&buf, binary.LittleEndian, Header{
		Head:           255,
		Version:        0,
		Session:        c.Session,
		SequenceNumber: sequence,
		MsgID:          msgID,
		Len:            dataLen,
	}); err != nil {
		return err
	}

	if data != nil {
		// write the data
		err := binary.Write(&buf, binary.LittleEndian, data)
		if err != nil {
			return err
		}

		// add magic bytes at the end
		err = binary.Write(&buf, binary.LittleEndian, magicEnd)
		if err != nil {
			return err
		}
	}

	// send all the data to the device
	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Write writes something to the device, wow!
func (c *Client) Write(msgID uint16, data []byte, dataLen uint32) error {
	if err := c.write(msgID, data, dataLen, c.Sequence); err != nil {
		return err
	}

	c.Sequence++
	return nil
}

// Read attempts to read a header and then the bytes.
func (c *Client) Read(excludeMagic bool) (*Header, []byte, error) {
	var p Header
	var b = make([]byte, 20)

	err := c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		return nil, nil, err
	}

	_, err = c.conn.Read(b)
	if err != nil {
		return nil, nil, err
	}

	err = binary.Read(bytes.NewReader(b), binary.LittleEndian, &p)
	if err != nil {
		return nil, nil, err
	}

	c.Sequence++

	if p.Len <= 0 || p.Len >= 100000 {
		return nil, nil, fmt.Errorf("invalid bodylength: %v", p.Len)
	}

	body := make([]byte, p.Len)
	err = binary.Read(c.conn, binary.LittleEndian, &body)
	if err != nil {
		return nil, nil, err
	}

	if excludeMagic && len(body) > 2 && bytes.Compare(body[len(body)-2:], []byte{10, 0}) == 0 {
		body = body[:len(body)-2]
	}

	return &p, body, nil
}

func (c *Client) Instruct(opcode uint16, name string, data map[string]interface{}) (*Header, []byte, error) {
	if len(name) < 1 {
		name = mappedOpcodes[opcode]
	}

	if data == nil {
		data = map[string]interface{}{"Name": name, "SessionID": fmt.Sprintf("0x%08X", c.Session)}
	} else {
		data["Name"] = name
		data["SessionID"] = fmt.Sprintf("0x%08X", c.Session)
	}

	params, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	err = c.Write(opcode, params, uint32(len(params))+2)
	if err != nil {
		return nil, nil, err
	}

	return c.Read(true)
}

func (c *Client) InstructRaw(opcode uint16, sequence uint32, data []byte, dataLen uint32) (*Header, []byte, error) {
	err := c.write(opcode, data, dataLen, sequence)
	if err != nil {
		return nil, nil, err
	}

	return c.Read(true)
}
