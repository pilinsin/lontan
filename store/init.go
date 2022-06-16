package store

import (
	crdt "github.com/pilinsin/p2p-verse/crdt"
	isign "github.com/pilinsin/util/crypto/sign"
	ed25519 "github.com/pilinsin/util/crypto/sign/ed25519"
)

type ISignKey = isign.ISignKey
type IVerfKey = isign.IVerfKey

var NewKeyPair = ed25519.NewKeyPair
var UnmarshalSign = ed25519.UnmarshalSignKey
var UnmarshalVerf = ed25519.UnmarshalVerfKey

func genKp() (crdt.IPrivKey, crdt.IPubKey, error) {
	kp := NewKeyPair()
	return kp.Sign(), kp.Verify(), nil
}
func marshalPub(pub crdt.IPubKey) ([]byte, error) {
	vk := pub.(isign.IVerfKey)
	return vk.Raw()
}
func unmarshalPub(m []byte) (crdt.IPubKey, error) {
	return UnmarshalVerf(m)
}

func init() {
	crdt.InitCryptoFuncs(genKp, marshalPub, unmarshalPub)
}
