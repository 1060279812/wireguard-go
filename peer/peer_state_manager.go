package peer

import "fmt"

// PeerStateManager 用来管理和触发状态变化
type PeerStateManager struct {
    listener PeerStateChangeListener
}

// NewPeerStateManager 创建一个新的 PeerStateManager
func NewPeerStateManager(listener PeerStateChangeListener) *PeerStateManager {
    return &PeerStateManager{listener: listener}
}

// TriggerStateChange 根据状态变化触发回调
func (psm *PeerStateManager) TriggerStateChange(state State) {
    switch state {
    case Starting:
        psm.listener.OnStarting()
    case Connecting:
        psm.listener.OnConnecting()
    case HandshakeFailedForOther:
        psm.listener.OnHandshakeFailedForOther()
    case HandshakeFailedForNetwork:
        psm.listener.OnHandshakeFailedForNetwork()
    case Disconnecting:
        psm.listener.OnDisconnecting()
    case Stopping:
        psm.listener.OnStopping()
    default:
        fmt.Println("Unknown state")
    }
}
