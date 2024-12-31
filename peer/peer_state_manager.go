package peerState

import (
	"sync"
)

// PeerStateManager 单例模式管理所有监听器
/**
sudo -i
su - galen
*/
type PeerStateManager struct {
	listener                 PeerStateChangeListener
	peerHandshakeFailedCount map[string]int // 记录握手失败次数 == 2次认为断开连接
	lastState                State          // 记录上一次状态
	lock                     sync.Mutex
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
func (manager *PeerStateManager) SetListener(listener PeerStateChangeListener) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.listener = listener

	// manager.listeners = append(manager.listeners, listener)
}

// RemoveListener 从管理器中移除指定的监听器
// func (manager *PeerStateManager) RemoveListener(listener PeerStateChangeListener) {
// 	manager.lock.Lock()
// 	defer manager.lock.Unlock()

//		for i, l := range manager.listeners {
//			if l == listener {
//				manager.listeners = append(manager.listeners[:i], manager.listeners[i+1:]...)
//				fmt.Printf("Listener %v removed/n", listener)
//				break
//			}
//		}
//	}
//
// NotifyStateChange 通知所有监听器状态变化
func (manager *PeerStateManager) NotifyStateChange(publicKey [NoisePublicKeySize]byte, state State) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	//if state == Created || state == Starting || state == Stopping {
	//	//过滤重复状态回调
	//	return
	//}
	//if state == manager.lastState {
	//	//过滤重复状态回调
	//	return
	//}
	//if manager.lastState == HandshakeFailedForNetwork && state == HandshakeFailedForOther {
	//	// 过滤掉网络断开带来的HandshakeFailedForNetwork和HandshakeFailedForOther的重复回调
	//	manager.lastState = state
	//	return
	//}
	//manager.lastState = state

	//if state == HandshakeSuccess || state == HandshakeFailedForOther || state == HandshakeFailedForNetwork {
	////	// 锁定当前 goroutine 到操作系统线程
	//runtime.LockOSThread()
	//defer runtime.UnlockOSThread()
	//}

	//java.OnMainLoad(func(reg java.Register) {
	//	reg.WithClass("com.nk.Hello").
	//		BindNative("nice", "void(java.lang.String[])", nice).
	//		Done()
	//})
	//

	//vm.RunSource("java.lang.Thread.currentThread[java.lang.Thread()]();").
	//	AsObject().
	//	Invoke("getName", "java.lang.String()").
	//	AsString()

	//com/wireguard/android/backend/GoBackend
	//vm.RunSource("java.lang.Thread.currentThread[com.wireguard.android.backend.GoBackend()]();").
	//	AsObject().
	//	Invoke("onPeerStateChange", "([BI)V", cPublicKeyC, C.int(state))

	//native.LoadClass("com.wireguard.android.backend.GoBackend()").New().Invoke("textGoCall", "void()")

	//vm.RunSource("java.lang.Thread.currentThread[com.wireguard.android.backend.GoBackend()]();").
	//	AsObject().
	//	Invoke("textGoCall", "")

	//cPublicKey := unsafe.Pointer(&publicKey[0]) // 获取数组的第一个元素的指针，转换为 unsafe.Pointer
	//// 将 unsafe.Pointer 转换为 C.uchar* 类型
	//cPublicKeyC := (*C.uchar)(cPublicKey) // 强制转换为 C 中的 uchar 类型指针
	//// 调用JNI函数
	//C.CallJavaOnPeerStateChange(cPublicKeyC, C.int(state))

	//doJNIOperation(cPublicKeyC, state)

	// Destroy JVM when done
	// if jvm != nil {
	//    C.JNI_DestroyJavaVM(jvm)
	// }

	// if(len(manager.listeners) == 0 || state == manager.lastState) {
	// 	//过滤重复状态回调
	// 	return
	// }
	// if(manager.lastState == HandshakeFailedForNetwork && state == HandshakeFailedForOther) {
	// 	// 过滤掉网络断开带来的HandshakeFailedForNetwork和HandshakeFailedForOther的重复回调
	// 	manager.lastState = state
	// 	return
	// }
	// manager.lastState = state
	// for _, listener := range manager.listeners {
	// 	listener.OnStateChange(publicKey,state)
	// }
}

// DestroyAllListeners 销毁所有监听器
func (manager *PeerStateManager) Destroy() {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	//once.Do(func() {
	if manager.listener != nil {
		manager.listener.Destroy()
	}
	//})

	// for _, listener := range manager.listeners {
	// 	listener.Destroy()
	// }
	// manager.listeners = nil // 清空监听器列表
}
