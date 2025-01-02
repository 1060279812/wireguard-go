package call

import peerState "github.com/1060279812/wireguard-go/peer"

type Platform int8

const (
	Android Platform = 1 //"android"
	IOS     Platform = 2 // "ios"
	MAC     Platform = 3 //"mac"
)

// 定义一个全局的Channel
//var PeerStateCallChan = make(chan PeerStateCallStruct)

func NotifyPeerStateChange(peerStateS PeerStateCallStruct) {
	//PeerStateCallChan <- peerStateS
}

func closePeerStateCallChan() {
	//close(PeerStateCallChan)
}

type PeerStateCallStruct struct {
	Platform  Platform
	PublicKey [32]byte
	State     peerState.State
}
