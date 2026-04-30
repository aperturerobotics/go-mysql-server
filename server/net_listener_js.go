//go:build js

package server

import (
	"fmt"
	"net"
)

func newNetListener(protocol, address string) (net.Listener, error) {
	return nil, fmt.Errorf("net listener is unsupported in js builds")
}
