package peerState

/*
#include <jni.h>
#include <android/log.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

static pthread_mutex_t jni_mutex = PTHREAD_MUTEX_INITIALIZER;	// 用于保护 JNI 调用的互斥锁
JavaVM *g_jvm = NULL;  // 用于保存 JavaVM 引用

JNIEXPORT jint JNI_OnLoad(JavaVM *vm, void *reserved) {
    g_jvm = vm;  // 保存 JavaVM 引用
    return JNI_VERSION_1_6;
}

// 获取当前线程的JNIEnv
JNIEnv* getJNIEnv() {
    JNIEnv *env;
    // 从JavaVM获取JNIEnv，如果当前线程没有附加到JVM，会自动附加
    if ((*g_jvm)->GetEnv(g_jvm, (void**)&env, JNI_VERSION_1_6) != JNI_OK) {
        // 当前线程没有附加到JVM，手动附加
        if ((*g_jvm)->AttachCurrentThread(g_jvm, &env, NULL) != 0) {
            return NULL;
        }
    }
    return env;
}

void detachCurrentThread(){
    if(g_jvm != NULL){
		(*g_jvm)->DetachCurrentThread(g_jvm);
	}
}

void CallJavaOnPeerStateChange(unsigned char *publicKey, int state){

    JNIEnv* env = getJNIEnv();

    if(env == NULL){
        return;
    }
    // 获取 GoBackend 类
    jclass clazz = (*env)->FindClass(env, "com/wireguard/android/backend/GoBackend");
    if (clazz == NULL) {
        return;
    }

    // 获取 onPeerStateChange 方法的 ID
    jmethodID methodID = (*env)->GetStaticMethodID(env, clazz, "onPeerStateChange", "([BI)V");
    if (methodID == NULL) {
        return;
    }

    size_t length = strlen((char *)publicKey);

    // 创建 jbyteArray 来保存 unsigned char* 数据
    jbyteArray byteArray = (*env)->NewByteArray(env, length);

    // 将 unsigned char* 数据复制到 jbyteArray 中
    (*env)->SetByteArrayRegion(env, byteArray, 0, length, (const jbyte *)publicKey);

    // 调用 Java 层的 OnPeerStateChange 方法
    (*env)->CallStaticVoidMethod(env,clazz, methodID, byteArray,state);
//    // 释放 publicKey 字符串对象
//    (*env)->DeleteLocalRef(env, publicKey);
    // 删除局部引用
    (*env)->DeleteLocalRef(env, byteArray);

    // 释放内存
//    free(nullTerminated);
}

// 声明将要在 Go 代码中使用的 C 函数
//void CallJavaOnPeerStateChange(char *publicKey,int state){
//
//    pthread_mutex_lock(&jni_mutex);
//
//	__android_log_print(ANDROID_LOG_ERROR, "JNI", "-----CallJavaOnPeerStateChange() publicKey=%s,state=%d", publicKey,state);
//
//    JNIEnv *env = getJNIEnv();
//
//     if(env == NULL){
//        pthread_mutex_unlock(&jni_mutex);
//        return;
//    }
//
//
//    // 获取 GoBackend 类
//    jclass clazz = (*env)->FindClass(env, "com/wireguard/android/backend/GoBackend");
//    if (clazz == NULL) {
//	    // detachCurrentThread();
//	    pthread_mutex_unlock(&jni_mutex);
//        return;
//    }
//
//    // 获取 onPeerStateChange 方法的 ID
//    jmethodID methodID = (*env)->GetStaticMethodID(env, clazz, "onPeerStateChange", "(Ljava/lang/String;I)V");
//    if (methodID == NULL) {
//	    // detachCurrentThread();
//	     pthread_mutex_unlock(&jni_mutex);
//        return;
//    }
//
//    // 使用 NewStringUTF 方法将 C 字符串转换为 Java 字符串 (jstring)
//    jstring publicKeyStr = (*env)->NewStringUTF(env, publicKey);
//    // 调用 Java 层的 OnPeerStateChange 方法
//    (*env)->CallStaticVoidMethod(env,clazz, methodID, publicKeyStr,state);
////    // 释放 publicKey 字符串对象
////    (*env)->DeleteLocalRef(env, publicKey);
//
//    //  detachCurrentThread();
//
//     pthread_mutex_unlock(&jni_mutex);
// }

// Helper function to initialize JVM and get JNIEnv
JNIEnv* createJNIEnv(JavaVM **jvm) {
    JavaVMInitArgs vmArgs;
    JavaVMOption options[1];

    // Set classpath to find Java classes
    options[0].optionString = "-Djava.class.path=.";
    vmArgs.version = JNI_VERSION_1_6;
    vmArgs.nOptions = 1;
    vmArgs.options = options;
    vmArgs.ignoreUnrecognized = JNI_FALSE;

    JNIEnv *env;
    // int res = JNI_CreateJavaVM(jvm, (void**)&env, &vmArgs);
    // if (res < 0) {
    //     return NULL;
    // }
    return env;
}
*/
import "C"
import (
	"sync"
	"unsafe"
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

	//if state == manager.lastState {
	//	//过滤重复状态回调
	//	return
	//}
	//if manager.lastState == HandshakeFailedForNetwork && state == HandshakeFailedForOther {
	//	// 过滤掉网络断开带来的HandshakeFailedForNetwork和HandshakeFailedForOther的重复回调
	//	manager.lastState = state
	//	return
	//}

	// 将[32]byte转换为字符串
	//var publicKeyStr = string(publicKey[:])
	// 使用C.CString将Go字符串转换为C字符串
	//publicKeyStr := string(publicKey[:])
	//defer C.free(unsafe.Pointer(cPublicKey))

	// 将 [32]byte 转换为 []byte
	//byteSlice := publicKey[:]
	//cPublicKey := C.CBytes(byteSlice)

	// 将 publicKey 转换为 unsafe.Pointer 类型
	cPublicKey := unsafe.Pointer(&publicKey[0]) // 获取数组的第一个元素的指针，转换为 unsafe.Pointer

	// 将 unsafe.Pointer 转换为 C.uchar* 类型
	cPublicKeyC := (*C.uchar)(cPublicKey) // 强制转换为 C 中的 uchar 类型指针

	// 注意：你在传递给 C 代码后，记得在 C 代码中释放内存
	//defer C.free(cPublicKey)

	// Initialize JVM and obtain JNIEnv
	//   var jvm *C.JavaVM
	//   env := C.createJNIEnv(&jvm)
	//   if env == nil {
	//      fmt.Println("Failed to create JNIEnv")
	//      return
	//  }
	// 调用JNI函数
	C.CallJavaOnPeerStateChange(cPublicKeyC, C.int(state))

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
