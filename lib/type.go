package lib

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

type IPAddr struct {
	IP    net.IP
	IPNet *net.IPNet
}

func NewIPAddr(addr string) (*IPAddr, error) {
	a := new(IPAddr)
	err := a.UnmarshalText([]byte(addr))
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a IPAddr) String() string {
	b, _ := a.MarshalText()
	return string(b)
}

func (a IPAddr) MarshalText() (text []byte, err error) {
	text, err = a.IP.MarshalText()
	if err != nil {
		return nil, err
	}
	if a.IPNet != nil {
		ones, _ := a.IPNet.Mask.Size()
		text = append(text, []byte(fmt.Sprintf("/%d", ones))...)
	}
	return text, nil
}

func (a *IPAddr) UnmarshalText(text []byte) error {
	if bytes.Index(text, []byte{'/'}) == -1 {
		ip := net.ParseIP(string(text))
		if ip == nil {
			return errors.New("Invalid IP address text was passed")
		}
		a.IP = ip
	} else {
		ip, ipnet, err := net.ParseCIDR(string(text))
		if err != nil {
			return err
		}
		a.IP = ip
		a.IPNet = ipnet
	}
	return nil
}

type IPAddrList []IPAddr

func (as IPAddrList) String() string {
	b, _ := as.MarshalText()
	return string(b)
}

func (as IPAddrList) MarshalText() (text []byte, err error) {
	var tmp [][]byte
	for _, a := range as {
		txt, err := a.MarshalText()
		if err != nil {
			return nil, err
		}
		tmp = append(tmp, txt)
	}
	return bytes.Join(tmp, []byte{' '}), nil
}

type CustomNs int

func (ns CustomNs) MarshalText() (text []byte, err error) {
	if ns == CustomNs(0) {
		return []byte{'0'}, nil
	}
	return []byte{'1'}, nil
}

type Timestamp struct {
	time.Time
}

const (
	DataTimestampFormat string = "2006-01-02 15:04:05.000000-0700"
	ArgTimestampFormat  string = "2006-01-02 15:04 MST"
)

func (t Timestamp) String() string {
	b, _ := t.MarshalText()
	return string(b)
}

func (t Timestamp) MarshalText() ([]byte, error) {
	return []byte(t.Format(DataTimestampFormat)), nil
}

func parseDateTimestampFormat(text []byte) (time.Time, error) {
	var sep byte
	if bytes.IndexRune(text, '+') != -1 {
		sep = '+'
	} else if bytes.IndexRune(text, '-') != 1 {
		sep = '-'
	} else {
		return time.Time{}, errors.New("Invalid timestamp")
	}
	dtzone := bytes.Split(text, []byte{sep})
	if len(dtzone) != 2 {
		return time.Time{}, errors.New("Invalid timestamp")
	}

	var norm []byte
	norm = append(norm, dtzone[0]...)

	n := 6
	if bytes.IndexRune(dtzone[0], '.') != -1 {
		dtms := bytes.Split(dtzone[0], []byte{'.'})
		if len(dtms) != 2 {
			return time.Time{}, errors.New("Invalid timestamp")
		}
		n -= len(dtms[1])
	} else {
		norm = append(norm, '.')
	}
	norm = append(norm, bytes.Repeat([]byte{'0'}, n)...)
	norm = append(norm, sep)

	n = 4
	n -= len(dtzone[1])
	norm = append(norm, dtzone[1]...)
	norm = append(norm, bytes.Repeat([]byte{'0'}, n)...)

	return time.Parse(DataTimestampFormat, string(norm))
}

func (t *Timestamp) UnmarshalText(text []byte) error {
	if tm, err := parseDateTimestampFormat(text); err == nil {
		*t = Timestamp{tm}
		return nil
	}
	if tm, err := time.Parse(time.UnixDate, string(text)); err == nil {
		*t = Timestamp{tm}
		return nil
	}
	return errors.New("Can't unmarshal timestamp: " + string(text))
}
