package rotserver

import (
	"log"
	"net"
	"encoding/hex"
	"rotorctl/internal/rotors"
)

type RotServer struct {
	rotor rotors.Rotor
	listenAddr string
	ln net.Listener
	quitCh chan struct{}
}


func NewRotServer(listenAddr string, rotor rotors.Rotor) *RotServer {
	return &RotServer {
		rotor: rotor,
		listenAddr: listenAddr,
		quitCh: make(chan struct{}),	
	}
}

func (s *RotServer) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	log.Println("👂 Listening on:", s.listenAddr)


	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitCh

	return nil
}

func (s *RotServer) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("❌ Accept error:", err)
			continue
		}
		
		log.Println("🔌 Incoming connection from:", conn.RemoteAddr())
		
		go s.readLoop(conn)
	}
}

func (s *RotServer) readLoop(conn net.Conn) {
	defer conn.Close()
	
	buf := make([]byte, 128)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("❌ Read error:", err)
			if err.Error() == "EOF" {
				return
			} else {
				continue
			}
		}

		// check for EOT (ctrl-d)
		if buf[0] == 0x04 {
			log.Println("🛑 Closing connection to:", conn.RemoteAddr())
			return
		}
		
		msg := buf[:n]
		log.Println("📥 Incoming:", hex.EncodeToString(msg))

		packet := RotServerData{}
		err = packet.Parse(buf[:n])
		if err != nil {
			log.Println("❌ Error parsing packet:", err)
			continue
		}
		log.Printf("📦 Cmd: %s, Az: %v, El: %v, Flags: %s\n", packet.Cmd, packet.Az, packet.El, packet.Flags)

		if packet.Cmd == CommandGet {
			az, el, err := s.rotor.GetPos()
			if err != nil {
				log.Println("❌ Error getting data from rotor:", err)
				continue
			} 
			log.Printf("⚙️  Rotor, Az: %v, El: %v\n", az, el)
			packet = RotServerData{Cmd: CommandStatus, Az: az, El: el, Flags: "OK"}
			data, _ := packet.toBytes()
			n, err := conn.Write(data)
			if err != nil {
				log.Println("❌ Error sending data to client:", n, err)
			}
		}

		if packet.Cmd == CommandSet {
			flags := "OK"
			
			err := s.rotor.SetPos(packet.Az, packet.El)
			if err != nil {
				log.Println("❌ Error sending data to rotor:", err)
				flags = "ERR"
			}
			
			packet = RotServerData{Cmd: CommandStatus, Az: packet.Az, El: packet.El, Flags: flags}
			data, _ := packet.toBytes()
			n, err := conn.Write(data)
			if err != nil {
				log.Println("❌ Error sending data to client:", n, err)
			}
		}
		
	}
}
