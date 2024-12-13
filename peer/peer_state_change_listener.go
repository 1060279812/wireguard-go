package peer

// 定义状态常量
type State int

const (
    Starting State = iota
    Connecting
    HandshakeFailedForOther
    HandshakeFailedForNetwork
    Disconnecting
    Stopping
)

// PeerStateChangeListener 接口定义了各个状态变化的回调方法
type PeerStateChangeListener interface {
    OnStarting()
    OnConnecting()
    OnHandshakeFailedForOther()
    OnHandshakeFailedForNetwork()
    OnDisconnecting()
    OnStopping()
}