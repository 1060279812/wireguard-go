package peerState

// 定义状态常量
type State int8

const (
	NoisePublicKeySize = 32
)

const (
	Created                   State = 1 //"created"
	Starting                  State = 2 // "starting"
	HandshakeSuccess          State = 3 //"handshakeSuccess"
	HandshakeFailedForOther   State = 4 //"handshakeFailedForOther"
	HandshakeFailedForNetwork State = 5 //"handshakeFailedForNetwork"
	Stopping                  State = 6 //"stopping"
)

// PeerStateChangeListener 定义回调接口
type PeerStateChangeListener interface {
	OnStateChange(noisePublicKey [NoisePublicKeySize]byte, state State) // 状态变化时回调
	Destroy()                                                           // 销毁回调
}
