package rot2prog

import (
	//"encoding/hex"
	"errors"
	"log"
	//"time"
	"github.com/albenik/go-serial/v2"
)


type Pos struct {
	Az float64
	El float64
}

type Rot2Prog struct {
	Device string
	CurPos Pos
	Port *serial.Port
}

func (r *Rot2Prog) Init() error {
	var err error
	r.Port, err = serial.Open(
		r.Device,
		serial.WithBaudrate(600),
		serial.WithDataBits(8),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(1000),
		serial.WithWriteTimeout(1000),
	)
    if err != nil {
		return err
    }
	if r.Port == nil {
		return errors.New("Device is NIL")
	}

	// test device
	s := Rot2ProgPacket { Cmd: Status }
	_ = s.GenBytes()
	n, err := r.Port.Write(s.Bytes)
    if err != nil {
		return err
    }
	if n < CmdPacketBytes {
		return errors.New("Can't send whole packet")
	}
	
	buf := make([]byte, StatusPacketBytes)
	n, err = r.Port.Read(buf)
    if err != nil {
		return err 
    }
	if n < StatusPacketBytes {
		log.Printf("Read: %v byte(s)\n", n) 
		return errors.New("Packet too short")
	}

	err = s.ParseBytes(buf)
	if err != nil {
		return err
	}
	
	log.Println("🔗 Connected to: ", r.Port)
	
	return nil
}

func (r *Rot2Prog) GetPos() (float64, float64, error) {
	// send status comand packet
	s := Rot2ProgPacket { Cmd: Status }
	_ = s.GenBytes()
	n, err := r.Port.Write(s.Bytes)
    if err != nil {
		return 0, 0, err
    }
	if n < CmdPacketBytes {
		return 0, 0, errors.New("Can't send whole packet")
	}
	// read answer
	buf := make([]byte, StatusPacketBytes)
	n, err = r.Port.Read(buf)
    if err != nil {
		return 0, 0, err 
    }
	if n < StatusPacketBytes {
		return 0, 0, errors.New("Didn't receive whole packet")
	}
	err = s.ParseBytes(buf)
	if err != nil {
		return 0, 0, err
	}

	r.CurPos.Az, r.CurPos.El = s.Az, s.El
	return s.Az, s.El, nil	
}


func (r *Rot2Prog) SetPos(az float64, el float64) error {
	// send status comand packet
	s := Rot2ProgPacket { Cmd: Set,  Az: az, El: el, DegPulse: int(DegPulse05)}
	_ = s.GenBytes()
	n, err := r.Port.Write(s.Bytes)
    if err != nil {
		return err
    }
	if n < CmdPacketBytes {
		return errors.New("Can't send whole packet")
	}
	// read answer
	buf := make([]byte, StatusPacketBytes)
	n, err = r.Port.Read(buf)
    if err != nil {
		return err 
    }
	if n < StatusPacketBytes {
		return errors.New("Didn't receive whole packet")
	}
	err = s.ParseBytes(buf)
	if err != nil {
		return err
	}
	
	return nil	
}

