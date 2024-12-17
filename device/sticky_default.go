//go:build !linux

package device

import (
	"github.com/1060279812/wireguard-go/conn"
	"github.com/1060279812/wireguard-go/rwcancel"
)

func (device *Device) startRouteListener(bind conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
