package peerState

// 定义状态常量
type State string

const (
	NoisePublicKeySize    = 32
)

const (
	Created                   State = "created"
	Starting                  State = "starting"
	// Connecting                State = "connecting"
    HandshakeSuccess          State = "handshakeSuccess"
	HandshakeFailedForOther   State = "handshakeFailedForOther"
	HandshakeFailedForNetwork State = "handshakeFailedForNetwork"
	// Disconnecting             State = "disconnecting"
	Stopping                  State = "stopping"
)

// PeerStateChangeListener 定义回调接口
type PeerStateChangeListener interface {
	OnStateChange(noisePublicKey [NoisePublicKeySize]byte,intstate State) // 状态变化时回调
	Destroy()                  // 销毁回调
}
