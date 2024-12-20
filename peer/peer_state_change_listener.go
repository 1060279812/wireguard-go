package peerState

// 定义状态常量
type State int8

const (
	NoisePublicKeySize   = 32
)

const (
	Created                   State = 1 //"created"
	Starting                  State = 2 // "starting"
	// Connecting                State =  "connecting"
    HandshakeSuccess          State = 3 //"handshakeSuccess"
	HandshakeFailedForOther   State =4 //"handshakeFailedForOther"
	HandshakeFailedForNetwork State =5 //"handshakeFailedForNetwork"
	// Disconnecting             State = "disconnecting"
	Stopping                  State = 6 //"stopping"
)

// PeerStateChangeListener 定义回调接口
type PeerStateChangeListener interface {
	OnStateChange(noisePublicKey [NoisePublicKeySize]byte,state State) // 状态变化时回调
	Destroy()                  // 销毁回调
}
//android项目用go语言写了一个
// type PeerStateChangeListener interface {
// 	OnStateChange(noisePublicKey [NoisePublicKeySize]byte,intstate State) // 状态变化时回调
// 	Destroy()                  // 销毁回调
// } 
// 和
// func GetInstance() *PeerStateManager 
// 在go层回调各种状态最终将状态回调给java层
