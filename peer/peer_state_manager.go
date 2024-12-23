package peerState

/*

#include <stdio.h>

// C function to add two integers
int add(int a, int b) {
    return a + b;
}

*/

import "C"
import (
	"fmt"
	"sync"
)

// PeerStateManager 单例模式管理所有监听器
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

// 	for i, l := range manager.listeners {
// 		if l == listener {
// 			manager.listeners = append(manager.listeners[:i], manager.listeners[i+1:]...)
// 			fmt.Printf("Listener %v removed/n", listener)
// 			break
// 		}
// 	}
// }

// NotifyStateChange 通知所有监听器状态变化
func (manager *PeerStateManager) NotifyStateChange(publicKey [NoisePublicKeySize]byte, state State) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if state == manager.lastState {
		//过滤重复状态回调
		return
	}
	if manager.lastState == HandshakeFailedForNetwork && state == HandshakeFailedForOther {
		// 过滤掉网络断开带来的HandshakeFailedForNetwork和HandshakeFailedForOther的重复回调
		manager.lastState = state
		return
	}

	// C.onStateChange2()
	
	// Calling the C function 'add' from Go
    a, b := 3, 5
    sum := C.add(C.int(a), C.int(b))

    fmt.Printf("The sum of %d and %d is: %d\n", a, b, sum)

	// 将[32]byte转换为字符串
	// var publicKeyStr = string(publicKey[:])
	// // 使用C.CString将Go字符串转换为C字符串
	// cPublicKey := C.CString(publicKeyStr)
	// // defer C.free(unsafe.Pointer(cPublicKey))

	//  // Initialize JVM and obtain JNIEnv
	//  var jvm *C.JavaVM
	//  env := C.createJNIEnv(&jvm)
	//  if env == nil {
	//     fmt.Println("Failed to create JNIEnv")
	//     return
	// }
	// // 调用JNI函数
	// C.CallJavaOnPeerStateChange(env,C.jstring(cPublicKey), C.jint(state))

	manager.lastState = state

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

	if manager.listener != nil {
		manager.listener.Destroy()
	}

	// for _, listener := range manager.listeners {
	// 	listener.Destroy()
	// }
	// manager.listeners = nil // 清空监听器列表
}
