package peer

import "fmt"

// MyPeerListener 实现 PeerStateChangeListener 接口
type MyPeerListener struct{}

// OnStarting 状态回调
func (mpl *MyPeerListener) OnStarting() {
    fmt.Println("State changed to Starting")
}

// OnConnecting 状态回调
func (mpl *MyPeerListener) OnConnecting() {
    fmt.Println("State changed to Connecting")
}

// OnHandshakeFailedForOther 状态回调
func (mpl *MyPeerListener) OnHandshakeFailedForOther() {
    fmt.Println("State changed to HandshakeFailedForOther")
}

// OnHandshakeFailedForNetwork 状态回调
func (mpl *MyPeerListener) OnHandshakeFailedForNetwork() {
    fmt.Println("State changed to HandshakeFailedForNetwork")
}

// OnDisconnecting 状态回调
func (mpl *MyPeerListener) OnDisconnecting() {
    fmt.Println("State changed to Disconnecting")
}

// OnStopping 状态回调
func (mpl *MyPeerListener) OnStopping() {
    fmt.Println("State changed to Stopping")
}
