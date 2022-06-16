package store

import (
	"encoding/base64"
	"errors"
	proto "google.golang.org/protobuf/proto"

	pb "github.com/pilinsin/lontan/store/pb"
)

type UserIdentity struct {
	userName string
	verfKey  IVerfKey
	signKey  ISignKey
}

func NewUserIdentity(name string, verf IVerfKey, sign ISignKey) *UserIdentity {
	return &UserIdentity{name, verf, sign}
}
func (ui UserIdentity) UserName() string { return ui.userName }
func (ui UserIdentity) Verify() IVerfKey { return ui.verfKey }
func (ui UserIdentity) Sign() ISignKey   { return ui.signKey }

func (ui *UserIdentity) Marshal() []byte {
	mv, _ := ui.verfKey.Raw()
	ms, _ := ui.signKey.Raw()
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

	verfKey, err := UnmarshalVerf(mui.GetVerf())
	if err != nil {
		return err
	}
	signKey, err := UnmarshalSign(mui.GetSign())
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
