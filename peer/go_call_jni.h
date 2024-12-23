//
// Created by GL on 2024/12/20.
//

#ifndef WIREGUARD_ANDROID_GO_CALL_JNI_H
#define WIREGUARD_ANDROID_GO_CALL_JNI_H

#include <jni.h>

void CallJavaOnPeerStateChange(JNIEnv *env,jstring publicKey,jint state);

#endif //WIREGUARD_ANDROID_GO_CALL_JNI_H
