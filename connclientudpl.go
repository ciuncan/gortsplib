package gortsplib

import (
	"net"
	"strconv"

	"github.com/aler9/gortsplib/multibuffer"
)

type connClientUDPListener struct {
	pc              net.PacketConn
	remoteIp        net.IP
	remoteZone      string
	remotePort      int
	udpFrameReadBuf *multibuffer.MultiBuffer
}

func newConnClientUDPListener(conf ConnClientConf, port int) (*connClientUDPListener, error) {
	pc, err := conf.ListenPacket("udp", ":"+strconv.FormatInt(int64(port), 10))
	if err != nil {
		return nil, err
	}

	return &connClientUDPListener{
		pc:              pc,
		udpFrameReadBuf: multibuffer.New(conf.ReadBufferCount, 2048),
	}, nil
}

func (l *connClientUDPListener) close() {
	l.pc.Close()
}

func (l *connClientUDPListener) read() ([]byte, error) {
	for {
		buf := l.udpFrameReadBuf.Next()
		n, addr, err := l.pc.ReadFrom(buf)
		if err != nil {
			return nil, err
		}

		uaddr := addr.(*net.UDPAddr)

		if !l.remoteIp.Equal(uaddr.IP) || l.remotePort != uaddr.Port {
			continue
		}

		return buf[:n], nil
	}
}

func (l *connClientUDPListener) write(buf []byte) error {
	_, err := l.pc.WriteTo(buf, &net.UDPAddr{
		IP:   l.remoteIp,
		Zone: l.remoteZone,
		Port: l.remotePort,
	})
	return err
}
