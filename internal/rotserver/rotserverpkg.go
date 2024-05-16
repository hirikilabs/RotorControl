package rotserver

import (
	"errors"
	"strconv"
	"strings"
)

const (
	FieldNumber int = 5
	Command int = 0
	Azimuth int = 1
	Elevation int = 2
	Flags int = 3
	Checksum int = 4

	CommandSet string = "SET"
	CommandGet string = "GET"
	CommandStatus string = "STATUS"
)

type RotServerData struct {
	Cmd string
	Az float64
	El float64
	Flags string
}

func (d *RotServerData) Parse(data []byte) error {
	// convert bytes to string
	s := string(data[:])

	// split string and check fields
	fields := strings.Split(strings.TrimSpace(s), ",")
	if len(fields) != FieldNumber {
		return errors.New("Wrong number of fields")
	}

	if fields[Command] != CommandSet && fields[Command] != CommandGet {
		return errors.New("Wrong command")
	}
	d.Cmd = fields[Command]

	// parse data
	// checksum
	checksum, err := strconv.ParseInt(fields[Checksum], 10, 32)
	if err != nil {
		return err
	}

	var calculated int64 = 0
	for i := 0; i < len(fields) - 1; i++ {
		for j := 0; j < len(fields[i]); j++ {
			calculated += int64(fields[i][j])
		}
	}
	if calculated != checksum {
		return errors.New("Wrong checksum")
	}

	d.Az, err = strconv.ParseFloat(fields[Azimuth], 64)
	if err != nil {
		return err
	}
	d.El, err = strconv.ParseFloat(fields[Elevation], 64)
	if err != nil {
		return err
	}

	d.Flags  = fields[Flags]
	
	return nil
}

