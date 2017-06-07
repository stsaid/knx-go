// Copyright 2017 Ole Krüger.

package proto

import (
	"errors"
	"fmt"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/util"
)

// A TunnelReq asks a gateway to transmit data.
type TunnelReq struct {
	// Communication channel
	Channel uint8

	// Sequential number, used to track acknowledgements
	SeqNumber uint8

	// Data to be tunneled
	Payload cemi.Message
}

// Service returns the service identifiers for tunnel requests.
func (TunnelReq) Service() ServiceID {
	return TunnelReqService
}

// Size returns the packed size.
func (req *TunnelReq) Size() uint {
	return 4 + cemi.Size(req.Payload)
}

// Pack the structure into the buffer.
func (req *TunnelReq) Pack(buffer []byte) {
	buffer[0] = 4
	buffer[1] = req.Channel
	buffer[2] = req.SeqNumber
	buffer[3] = 0
	cemi.Pack(buffer[4:], req.Payload)
}

// Unpack initializes the structure by parsing the given data.
func (req *TunnelReq) Unpack(data []byte) (n uint, err error) {
	var length, reserved uint8

	if n, err = util.UnpackSome(
		data, &length, &req.Channel, &req.SeqNumber, &reserved,
	); err != nil {
		return
	}

	if length != 4 {
		return n, errors.New("Length header is not 4")
	}

	m, err := cemi.Unpack(data[n:], &req.Payload)
	n += m

	return
}

// A TunnelResStatus is the status in a tunnel response.
type TunnelResStatus uint8

const (
	// TunnelResOk indicates that everything went as expected.
	TunnelResOk TunnelResStatus = 0x00

	// TunnelResUnsupported indicates that the CEMI-encoded frame inside the tunnel request was not
	// understood or is not supported.
	TunnelResUnsupported TunnelResStatus = 0x29
)

func (status TunnelResStatus) String() string {
	switch status {
	case TunnelResOk:
		return "Ok"

	case TunnelResUnsupported:
		return "Unsupported"

	default:
		return fmt.Sprintf("%#x", uint8(status))
	}
}

// A TunnelRes is a response to a TunnelRequest. It acts as an acknowledgement.
type TunnelRes struct {
	// Communication channel
	Channel uint8

	// Identifies the request that is being acknowledged
	SeqNumber uint8

	// Status code, determines whether the tunneling succeeded or not
	Status TunnelResStatus
}

// Service returns the service identifier for tunnel responses.
func (TunnelRes) Service() ServiceID {
	return TunnelResService
}

// Size returns the packed size.
func (TunnelRes) Size() uint {
	return 4
}

// Pack the structure into the buffer.
func (res *TunnelRes) Pack(buffer []byte) {
	buffer[0] = 4
	buffer[1] = res.Channel
	buffer[2] = res.SeqNumber
	buffer[3] = byte(res.Status)
}

// Unpack initializes the structure by parsing the given data.
func (res *TunnelRes) Unpack(data []byte) (n uint, err error) {
	var length uint8

	n, err = util.UnpackSome(data, &length, &res.Channel, &res.SeqNumber, (*uint8)(&res.Status))
	if err != nil {
		return
	}

	if length != 4 {
		return n, errors.New("Length header is not 4")
	}

	return
}
