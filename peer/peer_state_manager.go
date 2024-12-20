package peerState

/*
#cgo LDFLAGS: -L. -lBridgingGo
#include "go_call_jni.h"
#include <jni.h>

// 声明将要在 Go 代码中使用的 C 函数
void CallJavaOnPeerStateChange(JNIEnv *env, jobject obj, jstring state);

void CallJavaOnPeerStateChange(JNIEnv *env,jstring publicKey,jint state){

    // 获取 GoBackend 类
    jclass clazz = (*env)->FindClass(env, "com/wireguard/android/backend/GoBackend");
    if (clazz == NULL) {
        return;
    }

    // 获取 onPeerStateChange 方法的 ID
    jmethodID methodID = (*env)->GetStaticMethodID(env, clazz, "onPeerStateChange", "(Ljava/lang/String;I)V");
    if (methodID == NULL) {
        return;
    }

    // 调用 Java 层的 OnPeerStateChange 方法
    (*env)->CallStaticVoidMethod(env,clazz, methodID, publicKey,state);
//    // 释放 publicKey 字符串对象
//    (*env)->DeleteLocalRef(env, publicKey);
}
*/


import "C"
//go:generate go tool cgo peer_state_manager.go
import (
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
// 			fmt.Printf("Listener %v removed\n", listener)
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

	// 将[32]byte转换为字符串
	var publicKeyStr = string(publicKey[:])
	// 使用C.CString将Go字符串转换为C字符串
	cPublicKey := C.CString(publicKeyStr)
	// defer C.free(unsafe.Pointer(cPublicKey))

	// 调用JNI函数
	C.CallJavaOnPeerStateChange(C.jstring(cPublicKey), C.jint(state))

	manager.lastState = state

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
