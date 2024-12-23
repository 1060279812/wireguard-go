
#include "go_call_jni.h"

void CallJavaOnPeerStateChange(JNIEnv *env,jstring publicKey,jint state){

//    JNIEnv* env = getJNIEnv();

    if(env == NULL){
        return;
    }
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


