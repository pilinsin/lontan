package store

import(
	"errors"
	"encoding/base64"
	proto "google.golang.org/protobuf/proto"

	pb "github.com/pilinsin/lontan/store/pb"
	"github.com/pilinsin/util/crypto"
)

type UserIdentity struct{
	userName	string
	verfKey		crypto.IVerfKey
	signKey		crypto.ISignKey
}
func NewUserIdentity(name string, verf crypto.IVerfKey, sign crypto.ISignKey) *UserIdentity{
	return &UserIdentity{name, verf, sign}
}
func (ui UserIdentity) UserName() string{return ui.userName}
func (ui UserIdentity) Verify() crypto.IVerfKey{return ui.verfKey}
func (ui UserIdentity) Sign() crypto.ISignKey{return ui.signKey}

func (ui *UserIdentity) Marshal() []byte{
	mv, _ := crypto.MarshalVerfKey(ui.verfKey)
	ms, _ := crypto.MarshalSignKey(ui.signKey)
	mui := &pb.Identity{
		Name: ui.userName,
		Verf: mv,
		Sign: ms,
	}
	m, _ := proto.Marshal(mui)
	return m
}
func (ui *UserIdentity) Unmarshal(m []byte) error {
	mui := &pb.Identity{}
	if err := proto.Unmarshal(m, mui); err != nil {
		return err
	}

	verfKey, err := crypto.UnmarshalVerfKey(mui.GetVerf())
	if err != nil {
		return err
	}
	signKey, err := crypto.UnmarshalSignKey(mui.GetSign())
	if err != nil {
		return err
	}

	ui.userName = mui.GetName()
	ui.verfKey = verfKey
	ui.signKey = signKey
	return nil
}

func (ui UserIdentity) ToString() string {
	return base64.URLEncoding.EncodeToString(ui.Marshal())
}
func (ui *UserIdentity) FromString(addr string) error {
	if addr == "" {
		return errors.New("invalid addr")
	}
	m, err := base64.URLEncoding.DecodeString(addr)
	if err != nil {
		return err
	}
	return ui.Unmarshal(m)
}
