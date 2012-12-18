package main

import (
	"bufio"
	// "fmt"
	"io"
	"net"
	"os"
	"runtime"
)

// Almost same with net/textproto/reader.go ReadLine
func ReadLine(r *bufio.Reader) (string, error) {
	var line []byte
	for {
		l, more, err := r.ReadLine()
		if err != nil {
			return "", err
		}

		if line == nil && !more {
			return string(l), nil
		}
		line = append(line, l...)
		if !more {
			break
		}
	}
	return string(line), nil
}

func IsDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

type notification chan byte

func newNotification() notification {
	// Notification channle has size 1, so sending a single one will not block
	return make(chan byte, 1)
}

func (n notification) notify() {
	n <- 1
}

func (n notification) hasNotified() bool {
	select {
	case <-n:
		return true
	default:
		return false
	}
	return false
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

// Get host IP address
func hostIP() (addrs []string, err error) {
	name, err := os.Hostname()
	if err != nil {
		errl.Printf("Error get host name: %v\n", err)
		return
	}

	addrs, err = net.LookupHost(name)
	if err != nil {
		errl.Printf("Error getting host IP address: %v\n", err)
		return
	}
	return
}

func trimLastDot(s string) string {
	if len(s) > 0 && s[len(s)-1] == '.' {
		return s[:len(s)-1]
	}
	return s
}

func expandTilde(path string) string {
	if len(path) > 0 && path[0] == '~' {
		return homeDir + path[1:]
	}
	return path
}

func hostIsIP(host string) bool {
	host, _ = splitHostPort(host)
	return net.ParseIP(host) != nil
}

func copyN(r io.Reader, w, contBuf io.Writer, n int, buf, pre, end []byte) (err error) {
	var nn int
	bufLen := len(buf)
	var b []byte
	for n != 0 {
		if pre != nil {
			if len(pre) >= bufLen {
				if _, err = w.Write(pre); err != nil {
					return
				}
				pre = nil
				continue
			}
			// append pre to buf
			copy(buf, pre)
			if len(pre)+n < bufLen {
				b = buf[len(pre) : len(pre)+n]
			} else {
				b = buf[len(pre):]
			}
		} else {
			if n < bufLen {
				b = buf[:n]
			} else {
				b = buf
			}
		}
		if nn, err = r.Read(b); err != nil {
			return
		}
		n -= nn
		if pre != nil {
			nn += len(pre)
			pre = nil
		}
		if n == 0 && end != nil && nn+len(end) <= bufLen {
			copy(buf[nn:], end)
			nn += len(end)
			end = nil
		}
		if contBuf != nil {
			contBuf.Write(buf[:nn])
		}
		if _, err = w.Write(buf[:nn]); err != nil {
			return
		}
	}
	if end != nil {
		if _, err = w.Write(end); err != nil {
			return
		}
	}
	return
}
