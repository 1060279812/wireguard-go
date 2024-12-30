package peerState

/*
#include <jni.h>
#include <android/log.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

static pthread_mutex_t jni_mutex = PTHREAD_MUTEX_INITIALIZER;	// 用于保护 JNI 调用的互斥锁
JavaVM *g_jvm = NULL;  // 用于保存 JavaVM 引用
pthread_t curr_thread_id;

JNIEXPORT jint JNI_OnLoad(JavaVM *vm, void *reserved) {
    g_jvm = vm;  // 保存 JavaVM 引用
    __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "---------------JNI_OnLoad--------------");
    return JNI_VERSION_1_6;
}

// 获取当前线程的JNIEnv
JNIEnv* getJNIEnv() {
    JNIEnv *env;
    if(g_jvm == NULL){
        __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "---------------g_jvm == NULL--------------");
        return NULL;
    }
    // 从JavaVM获取JNIEnv，如果当前线程没有附加到JVM，会自动附加
    if ((*g_jvm)->GetEnv(g_jvm, (void**)&env, JNI_VERSION_1_6) != JNI_OK) {
       // 当前线程没有附加到JVM，手动附加
       if ((*g_jvm)->AttachCurrentThread(g_jvm, &env, NULL) != JNI_OK) {
           return NULL;
       }
    }
    __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "---------------getJNIEnv--------------");
    return env;
}

void detachCurrentThread(){
    if(g_jvm != NULL){
		(*g_jvm)->DetachCurrentThread(g_jvm);
	}
}

void get_thread_id() {
    pthread_t thread_id = pthread_self();
    //if(curr_thread_id != NULL && thread_id != curr_thread_id){
    //   detachCurrentThread();
    //}
    //curr_thread_id = thread_id;
    __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "Current thread ID: %p\n",(void*) thread_id);
}

void CallJavaOnPeerStateChange(unsigned char *publicKey, int state){

    //pthread_mutex_lock(&jni_mutex);

    get_thread_id();

    JNIEnv* env = getJNIEnv();

    if(env == NULL){
        __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "Failed to attach current thread");
        //pthread_mutex_unlock(&jni_mutex);
        return;
    }
    __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "Current JNIEnv Address: %p\n", env);
    // 获取 GoBackend 类
    jclass clazz = (*env)->FindClass(env, "com/wireguard/android/backend/GoBackend");
    if (clazz == NULL) {
        __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "---------------clazz == NULL--------------");
        //pthread_mutex_unlock(&jni_mutex);
        return;
    }

    // 获取 onPeerStateChange 方法的 ID
    jmethodID methodID = (*env)->GetStaticMethodID(env, clazz, "onPeerStateChange", "([BI)V");
    if (methodID == NULL) {
        __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "---------------methodID == NULL--------------");
        //pthread_mutex_unlock(&jni_mutex);
        return;
    }

    size_t length = strlen((char *)publicKey);

    // 创建 jbyteArray 来保存 unsigned char* 数据
    jbyteArray byteArray = (*env)->NewByteArray(env, length);

    // 将 unsigned char* 数据复制到 jbyteArray 中
    (*env)->SetByteArrayRegion(env, byteArray, 0, length, (const jbyte *)publicKey);

     __android_log_print(ANDROID_LOG_ERROR, "GoJNI", "---------------CallStaticVoidMethod--------------");

    // 调用 Java 层的 OnPeerStateChange 方法
    (*env)->CallStaticVoidMethod(env,clazz, methodID, byteArray,state);
    // 删除局部引用
    (*env)->DeleteLocalRef(env, byteArray);

    detachCurrentThread();

    //pthread_mutex_unlock(&jni_mutex);
}
*/
import "C"
import (
	//"fmt"
	//"github.com/shangzebei/gojni/java"
	//"github.com/shangzebei/gojni/native"
	"github.com/shangzebei/gojni/vm"
	"sync"
	"unsafe"
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

	cPublicKey := unsafe.Pointer(&publicKey[0]) // 获取数组的第一个元素的指针，转换为 unsafe.Pointer
	// 将 unsafe.Pointer 转换为 C.uchar* 类型
	cPublicKeyC := (*C.uchar)(cPublicKey) // 强制转换为 C 中的 uchar 类型指针

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
	//fmt.Println(vm.RunSource("java.lang.Thread.currentThread[java.lang.Thread()]();").
	//	AsObject().
	//	Invoke("getName", "java.lang.String()").
	//	AsString())

	//com/wireguard/android/backend/GoBackend
	vm.RunSource("java.lang.Thread.currentThread[java.lang.Thread()]();").
		AsObject().
		Invoke("onPeerStateChange", "([BI)V")

	// 调用JNI函数
	C.CallJavaOnPeerStateChange(cPublicKeyC, C.int(state))

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

//var (
//	jvm      C.JavaVMAlias // JVM指针
//	jniEnv   C.JNIEnv      // JNI环境指针
//	jniMutex sync.Mutex    // 用于同步JNI调用的互斥锁
//)
//
//// 初始化JVM
//func initJVM() {
//	// ... JVM初始化代码 ...
//	jvm = C.JavaVMAlias(unsafe.Pointer(C.g_jvm))
//	// ... 获取JNIEnv ...
//	jniEnv = C.JNIEnv(unsafe.Pointer((*C.JNIEnv)(C.getJNIEnv())))
//}
//
//// JNI同步调用
//func callJNI(cPublicKeyC *C.uchar, state State) {
//	jniMutex.Lock()
//	defer jniMutex.Unlock()
//	// ... 调用JNI函数 ...
//	C.CallJavaOnPeerStateChange(cPublicKeyC, C.int(state))
//}
//
//// 子Goroutine中安全调用JNI
//func doJNIOperation(cPublicKeyC *C.uchar, state State) {
//	// 确保JVM已初始化
//	if jvm == nil {
//		initJVM()
//	}
//	// 调用JNI方法
//	callJNI(cPublicKeyC, C.int(state))
//}
