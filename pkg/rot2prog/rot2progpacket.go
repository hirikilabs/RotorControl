package rot2prog

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	Stop byte = 0x0F
	Status byte = 0x1F
	Set byte = 0x2F

	DegPulse1 byte = 0x01
	DegPulse05 byte = 0x02
	DegPulse02 byte = 0x04

	CmdPacketBytes int = 13
    StatusPacketBytes int = 12
)

type Rot2ProgPacket struct {
	Cmd byte
	Az float64
	El float64
	DegPulse int
	Bytes []byte
}

func (p* Rot2ProgPacket) GenBytes() error {
	p.Bytes = make([]byte, CmdPacketBytes)
	p.Bytes[0] = 0x57
	p.Bytes[12] = 0x20
    switch p.Cmd {
	case Set:
		p.Bytes[11] = Set
		// calculate pulses from coordinates
		azData := int64(float64(p.DegPulse) * (360.0 + p.Az))
		elData := int64(float64(p.DegPulse) * (360.0 + p.El))
		// convert to XXXX zero padded
		azString := fmt.Sprintf("%04d", azData)
		elString := fmt.Sprintf("%04d", elData)

		// fill bytes as numbers
		p.Bytes[1] = azString[0] 
		p.Bytes[2] = azString[1]
		p.Bytes[3] = azString[2]
		p.Bytes[4] = azString[3]

		p.Bytes[5] = DegPulse05

		p.Bytes[6] = elString[0]
		p.Bytes[7] = elString[1]
		p.Bytes[8] = elString[2]
		p.Bytes[9] = elString[3]

		p.Bytes[10] = DegPulse05
	case Stop:
		p.Bytes[11] = Stop
	case Status:
		p.Bytes[11] = Status
	default:
		return errors.New("No such command")
	}

	return nil
}

func (p* Rot2ProgPacket) ParseBytes(data []byte) error {
	if len(data) < StatusPacketBytes {
		return errors.New("Packet too short")
	}
	if data[0] != 0x57 || data[StatusPacketBytes-1] != 0x20 {
		return errors.New("Malformed packet")
	}

	// get Az and El
	azString := string(rune(data[1]+48)) + string(rune(data[2]+48)) + string(rune(data[3]+48)) + "." + string(rune(data[4]+48))
	elString := string(rune(data[6]+48)) + string(rune(data[7]+48)) + string(rune(data[8]+48)) + "." + string(rune(data[9]+48))

	var err error
	p.Az, err = strconv.ParseFloat(azString, 64)
	if err != nil {
		return err
	}
	p.Az = p.Az - 360
	
	p.El, err = strconv.ParseFloat(elString, 64)
	if err != nil {
		return err
	}
	p.El = p.El - 360

	if data[5] == data[10] {
		p.DegPulse = int(data[5])
	}
	
	return nil

}
