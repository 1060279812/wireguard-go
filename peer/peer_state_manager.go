package peerState

import "sync"
import "fmt"

// PeerStateManager 单例模式管理所有监听器
type PeerStateManager struct {
	listeners []PeerStateChangeListener
	peerHandshakeFailedCount map[string]int // 记录握手失败次数 == 2次认为断开连接
	lastState State  // 记录上一次状态
	lock      sync.Mutex
}

// Singleton instance
var instance *PeerStateManager
var once sync.Once

// GetInstance 获取单例实例
func GetInstance() *PeerStateManager {
	once.Do(func() {
		instance = &PeerStateManager{
			//peerHandshakeFailedCount:make(map[string]int),
			//lastState: Created,
		}
	})
	return instance
}

// AddListener 向管理器添加回调监听器
func (manager *PeerStateManager) AddListener(listener PeerStateChangeListener) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.listeners = append(manager.listeners, listener)
}

// RemoveListener 从管理器中移除指定的监听器
func (manager *PeerStateManager) RemoveListener(listener PeerStateChangeListener) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	for i, l := range manager.listeners {
		if l == listener {
			manager.listeners = append(manager.listeners[:i], manager.listeners[i+1:]...)
			fmt.Printf("Listener %v removed\n", listener)
			break
		}
	}
}

// NotifyStateChange 通知所有监听器状态变化
func (manager *PeerStateManager) NotifyStateChange(publicKey [NoisePublicKeySize]byte,state State) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if(len(manager.listeners) == 0 || state == manager.lastState) {
		return
	}
	if(manager.lastState == HandshakeFailedForNetwork && state == HandshakeFailedForOther) {
		// 过滤掉网络断开带来的HandshakeFailedForNetwork和HandshakeFailedForOther的重复回调
		manager.lastState = state
		return
	}
	manager.lastState = state
	for _, listener := range manager.listeners {
		listener.OnStateChange(publicKey,state)
	}
}

// DestroyAllListeners 销毁所有监听器
func (manager *PeerStateManager) DestroyAllListeners() {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	for _, listener := range manager.listeners {
		listener.Destroy()
	}
	manager.listeners = nil // 清空监听器列表
}
