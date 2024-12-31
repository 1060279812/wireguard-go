package call

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
import peerState "github.com/1060279812/wireguard-go/peer"

type Platform int8

const (
	Android Platform = 1 //"android"
	IOS     Platform = 2 // "ios"
	MAC     Platform = 3 //"mac"
)

// 定义一个全局的Channel
var GlobalPeerCallChannel = make(chan PeerStateCallStruct)

type PeerStateCallStruct struct {
	Platform  Platform
	PublicKey [32]byte
	State     peerState.State
}
