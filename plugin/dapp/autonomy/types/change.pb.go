// Code generated by protoc-gen-go. DO NOT EDIT.
// source: change.proto

package types

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AutonomyProposalChange struct {
	PropChange *ProposalChange `protobuf:"bytes,1,opt,name=propChange,proto3" json:"propChange,omitempty"`
	//
	CurRule *RuleConfig `protobuf:"bytes,2,opt,name=curRule,proto3" json:"curRule,omitempty"`
	//
	Board *ActiveBoard `protobuf:"bytes,3,opt,name=board,proto3" json:"board,omitempty"`
	//
	VoteResult *VoteResult `protobuf:"bytes,4,opt,name=voteResult,proto3" json:"voteResult,omitempty"`
	//
	Status               int32    `protobuf:"varint,5,opt,name=status,proto3" json:"status,omitempty"`
	Address              string   `protobuf:"bytes,6,opt,name=address,proto3" json:"address,omitempty"`
	Height               int64    `protobuf:"varint,7,opt,name=height,proto3" json:"height,omitempty"`
	Index                int32    `protobuf:"varint,8,opt,name=index,proto3" json:"index,omitempty"`
	ProposalID           string   `protobuf:"bytes,9,opt,name=proposalID,proto3" json:"proposalID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AutonomyProposalChange) Reset()         { *m = AutonomyProposalChange{} }
func (m *AutonomyProposalChange) String() string { return proto.CompactTextString(m) }
func (*AutonomyProposalChange) ProtoMessage()    {}
func (*AutonomyProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{0}
}

func (m *AutonomyProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AutonomyProposalChange.Unmarshal(m, b)
}
func (m *AutonomyProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AutonomyProposalChange.Marshal(b, m, deterministic)
}
func (m *AutonomyProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AutonomyProposalChange.Merge(m, src)
}
func (m *AutonomyProposalChange) XXX_Size() int {
	return xxx_messageInfo_AutonomyProposalChange.Size(m)
}
func (m *AutonomyProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_AutonomyProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_AutonomyProposalChange proto.InternalMessageInfo

func (m *AutonomyProposalChange) GetPropChange() *ProposalChange {
	if m != nil {
		return m.PropChange
	}
	return nil
}

func (m *AutonomyProposalChange) GetCurRule() *RuleConfig {
	if m != nil {
		return m.CurRule
	}
	return nil
}

func (m *AutonomyProposalChange) GetBoard() *ActiveBoard {
	if m != nil {
		return m.Board
	}
	return nil
}

func (m *AutonomyProposalChange) GetVoteResult() *VoteResult {
	if m != nil {
		return m.VoteResult
	}
	return nil
}

func (m *AutonomyProposalChange) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *AutonomyProposalChange) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *AutonomyProposalChange) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *AutonomyProposalChange) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *AutonomyProposalChange) GetProposalID() string {
	if m != nil {
		return m.ProposalID
	}
	return ""
}

// action
type ProposalChange struct {
	//
	Year  int32 `protobuf:"varint,1,opt,name=year,proto3" json:"year,omitempty"`
	Month int32 `protobuf:"varint,2,opt,name=month,proto3" json:"month,omitempty"`
	Day   int32 `protobuf:"varint,3,opt,name=day,proto3" json:"day,omitempty"`
	//
	Changes []*Change `protobuf:"bytes,4,rep,name=changes,proto3" json:"changes,omitempty"`
	//
	StartBlockHeight     int64    `protobuf:"varint,5,opt,name=startBlockHeight,proto3" json:"startBlockHeight,omitempty"`
	EndBlockHeight       int64    `protobuf:"varint,6,opt,name=endBlockHeight,proto3" json:"endBlockHeight,omitempty"`
	RealEndBlockHeight   int64    `protobuf:"varint,7,opt,name=realEndBlockHeight,proto3" json:"realEndBlockHeight,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProposalChange) Reset()         { *m = ProposalChange{} }
func (m *ProposalChange) String() string { return proto.CompactTextString(m) }
func (*ProposalChange) ProtoMessage()    {}
func (*ProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{1}
}

func (m *ProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProposalChange.Unmarshal(m, b)
}
func (m *ProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProposalChange.Marshal(b, m, deterministic)
}
func (m *ProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProposalChange.Merge(m, src)
}
func (m *ProposalChange) XXX_Size() int {
	return xxx_messageInfo_ProposalChange.Size(m)
}
func (m *ProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_ProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_ProposalChange proto.InternalMessageInfo

func (m *ProposalChange) GetYear() int32 {
	if m != nil {
		return m.Year
	}
	return 0
}

func (m *ProposalChange) GetMonth() int32 {
	if m != nil {
		return m.Month
	}
	return 0
}

func (m *ProposalChange) GetDay() int32 {
	if m != nil {
		return m.Day
	}
	return 0
}

func (m *ProposalChange) GetChanges() []*Change {
	if m != nil {
		return m.Changes
	}
	return nil
}

func (m *ProposalChange) GetStartBlockHeight() int64 {
	if m != nil {
		return m.StartBlockHeight
	}
	return 0
}

func (m *ProposalChange) GetEndBlockHeight() int64 {
	if m != nil {
		return m.EndBlockHeight
	}
	return 0
}

func (m *ProposalChange) GetRealEndBlockHeight() int64 {
	if m != nil {
		return m.RealEndBlockHeight
	}
	return 0
}

type Change struct {
	// 1    0
	Cancel               bool     `protobuf:"varint,1,opt,name=cancel,proto3" json:"cancel,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Change) Reset()         { *m = Change{} }
func (m *Change) String() string { return proto.CompactTextString(m) }
func (*Change) ProtoMessage()    {}
func (*Change) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{2}
}

func (m *Change) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Change.Unmarshal(m, b)
}
func (m *Change) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Change.Marshal(b, m, deterministic)
}
func (m *Change) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Change.Merge(m, src)
}
func (m *Change) XXX_Size() int {
	return xxx_messageInfo_Change.Size(m)
}
func (m *Change) XXX_DiscardUnknown() {
	xxx_messageInfo_Change.DiscardUnknown(m)
}

var xxx_messageInfo_Change proto.InternalMessageInfo

func (m *Change) GetCancel() bool {
	if m != nil {
		return m.Cancel
	}
	return false
}

func (m *Change) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type RevokeProposalChange struct {
	ProposalID           string   `protobuf:"bytes,1,opt,name=proposalID,proto3" json:"proposalID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RevokeProposalChange) Reset()         { *m = RevokeProposalChange{} }
func (m *RevokeProposalChange) String() string { return proto.CompactTextString(m) }
func (*RevokeProposalChange) ProtoMessage()    {}
func (*RevokeProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{3}
}

func (m *RevokeProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RevokeProposalChange.Unmarshal(m, b)
}
func (m *RevokeProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RevokeProposalChange.Marshal(b, m, deterministic)
}
func (m *RevokeProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RevokeProposalChange.Merge(m, src)
}
func (m *RevokeProposalChange) XXX_Size() int {
	return xxx_messageInfo_RevokeProposalChange.Size(m)
}
func (m *RevokeProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_RevokeProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_RevokeProposalChange proto.InternalMessageInfo

func (m *RevokeProposalChange) GetProposalID() string {
	if m != nil {
		return m.ProposalID
	}
	return ""
}

type VoteProposalChange struct {
	ProposalID           string   `protobuf:"bytes,1,opt,name=proposalID,proto3" json:"proposalID,omitempty"`
	Approve              bool     `protobuf:"varint,2,opt,name=approve,proto3" json:"approve,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoteProposalChange) Reset()         { *m = VoteProposalChange{} }
func (m *VoteProposalChange) String() string { return proto.CompactTextString(m) }
func (*VoteProposalChange) ProtoMessage()    {}
func (*VoteProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{4}
}

func (m *VoteProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoteProposalChange.Unmarshal(m, b)
}
func (m *VoteProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoteProposalChange.Marshal(b, m, deterministic)
}
func (m *VoteProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoteProposalChange.Merge(m, src)
}
func (m *VoteProposalChange) XXX_Size() int {
	return xxx_messageInfo_VoteProposalChange.Size(m)
}
func (m *VoteProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_VoteProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_VoteProposalChange proto.InternalMessageInfo

func (m *VoteProposalChange) GetProposalID() string {
	if m != nil {
		return m.ProposalID
	}
	return ""
}

func (m *VoteProposalChange) GetApprove() bool {
	if m != nil {
		return m.Approve
	}
	return false
}

type TerminateProposalChange struct {
	ProposalID           string   `protobuf:"bytes,1,opt,name=proposalID,proto3" json:"proposalID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TerminateProposalChange) Reset()         { *m = TerminateProposalChange{} }
func (m *TerminateProposalChange) String() string { return proto.CompactTextString(m) }
func (*TerminateProposalChange) ProtoMessage()    {}
func (*TerminateProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{5}
}

func (m *TerminateProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TerminateProposalChange.Unmarshal(m, b)
}
func (m *TerminateProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TerminateProposalChange.Marshal(b, m, deterministic)
}
func (m *TerminateProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TerminateProposalChange.Merge(m, src)
}
func (m *TerminateProposalChange) XXX_Size() int {
	return xxx_messageInfo_TerminateProposalChange.Size(m)
}
func (m *TerminateProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_TerminateProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_TerminateProposalChange proto.InternalMessageInfo

func (m *TerminateProposalChange) GetProposalID() string {
	if m != nil {
		return m.ProposalID
	}
	return ""
}

// receipt
type ReceiptProposalChange struct {
	Prev                 *AutonomyProposalChange `protobuf:"bytes,1,opt,name=prev,proto3" json:"prev,omitempty"`
	Current              *AutonomyProposalChange `protobuf:"bytes,2,opt,name=current,proto3" json:"current,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ReceiptProposalChange) Reset()         { *m = ReceiptProposalChange{} }
func (m *ReceiptProposalChange) String() string { return proto.CompactTextString(m) }
func (*ReceiptProposalChange) ProtoMessage()    {}
func (*ReceiptProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{6}
}

func (m *ReceiptProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReceiptProposalChange.Unmarshal(m, b)
}
func (m *ReceiptProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReceiptProposalChange.Marshal(b, m, deterministic)
}
func (m *ReceiptProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiptProposalChange.Merge(m, src)
}
func (m *ReceiptProposalChange) XXX_Size() int {
	return xxx_messageInfo_ReceiptProposalChange.Size(m)
}
func (m *ReceiptProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiptProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiptProposalChange proto.InternalMessageInfo

func (m *ReceiptProposalChange) GetPrev() *AutonomyProposalChange {
	if m != nil {
		return m.Prev
	}
	return nil
}

func (m *ReceiptProposalChange) GetCurrent() *AutonomyProposalChange {
	if m != nil {
		return m.Current
	}
	return nil
}

type LocalProposalChange struct {
	PropBd               *AutonomyProposalChange `protobuf:"bytes,1,opt,name=propBd,proto3" json:"propBd,omitempty"`
	Comments             []string                `protobuf:"bytes,2,rep,name=comments,proto3" json:"comments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *LocalProposalChange) Reset()         { *m = LocalProposalChange{} }
func (m *LocalProposalChange) String() string { return proto.CompactTextString(m) }
func (*LocalProposalChange) ProtoMessage()    {}
func (*LocalProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{7}
}

func (m *LocalProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LocalProposalChange.Unmarshal(m, b)
}
func (m *LocalProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LocalProposalChange.Marshal(b, m, deterministic)
}
func (m *LocalProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LocalProposalChange.Merge(m, src)
}
func (m *LocalProposalChange) XXX_Size() int {
	return xxx_messageInfo_LocalProposalChange.Size(m)
}
func (m *LocalProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_LocalProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_LocalProposalChange proto.InternalMessageInfo

func (m *LocalProposalChange) GetPropBd() *AutonomyProposalChange {
	if m != nil {
		return m.PropBd
	}
	return nil
}

func (m *LocalProposalChange) GetComments() []string {
	if m != nil {
		return m.Comments
	}
	return nil
}

// query
type ReqQueryProposalChange struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Count                int32    `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
	Direction            int32    `protobuf:"varint,4,opt,name=direction,proto3" json:"direction,omitempty"`
	Height               int64    `protobuf:"varint,5,opt,name=height,proto3" json:"height,omitempty"`
	Index                int32    `protobuf:"varint,6,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqQueryProposalChange) Reset()         { *m = ReqQueryProposalChange{} }
func (m *ReqQueryProposalChange) String() string { return proto.CompactTextString(m) }
func (*ReqQueryProposalChange) ProtoMessage()    {}
func (*ReqQueryProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{8}
}

func (m *ReqQueryProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqQueryProposalChange.Unmarshal(m, b)
}
func (m *ReqQueryProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqQueryProposalChange.Marshal(b, m, deterministic)
}
func (m *ReqQueryProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqQueryProposalChange.Merge(m, src)
}
func (m *ReqQueryProposalChange) XXX_Size() int {
	return xxx_messageInfo_ReqQueryProposalChange.Size(m)
}
func (m *ReqQueryProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqQueryProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_ReqQueryProposalChange proto.InternalMessageInfo

func (m *ReqQueryProposalChange) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *ReqQueryProposalChange) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *ReqQueryProposalChange) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *ReqQueryProposalChange) GetDirection() int32 {
	if m != nil {
		return m.Direction
	}
	return 0
}

func (m *ReqQueryProposalChange) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *ReqQueryProposalChange) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type ReplyQueryProposalChange struct {
	PropChanges          []*AutonomyProposalChange `protobuf:"bytes,1,rep,name=propChanges,proto3" json:"propChanges,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *ReplyQueryProposalChange) Reset()         { *m = ReplyQueryProposalChange{} }
func (m *ReplyQueryProposalChange) String() string { return proto.CompactTextString(m) }
func (*ReplyQueryProposalChange) ProtoMessage()    {}
func (*ReplyQueryProposalChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c013f0fbf0b6ffb, []int{9}
}

func (m *ReplyQueryProposalChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReplyQueryProposalChange.Unmarshal(m, b)
}
func (m *ReplyQueryProposalChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReplyQueryProposalChange.Marshal(b, m, deterministic)
}
func (m *ReplyQueryProposalChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReplyQueryProposalChange.Merge(m, src)
}
func (m *ReplyQueryProposalChange) XXX_Size() int {
	return xxx_messageInfo_ReplyQueryProposalChange.Size(m)
}
func (m *ReplyQueryProposalChange) XXX_DiscardUnknown() {
	xxx_messageInfo_ReplyQueryProposalChange.DiscardUnknown(m)
}

var xxx_messageInfo_ReplyQueryProposalChange proto.InternalMessageInfo

func (m *ReplyQueryProposalChange) GetPropChanges() []*AutonomyProposalChange {
	if m != nil {
		return m.PropChanges
	}
	return nil
}

func init() {
	proto.RegisterType((*AutonomyProposalChange)(nil), "types.AutonomyProposalChange")
	proto.RegisterType((*ProposalChange)(nil), "types.ProposalChange")
	proto.RegisterType((*Change)(nil), "types.Change")
	proto.RegisterType((*RevokeProposalChange)(nil), "types.RevokeProposalChange")
	proto.RegisterType((*VoteProposalChange)(nil), "types.VoteProposalChange")
	proto.RegisterType((*TerminateProposalChange)(nil), "types.TerminateProposalChange")
	proto.RegisterType((*ReceiptProposalChange)(nil), "types.ReceiptProposalChange")
	proto.RegisterType((*LocalProposalChange)(nil), "types.LocalProposalChange")
	proto.RegisterType((*ReqQueryProposalChange)(nil), "types.ReqQueryProposalChange")
	proto.RegisterType((*ReplyQueryProposalChange)(nil), "types.ReplyQueryProposalChange")
}

func init() {
	proto.RegisterFile("change.proto", fileDescriptor_4c013f0fbf0b6ffb)
}

var fileDescriptor_4c013f0fbf0b6ffb = []byte{
	// 596 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4d, 0x6f, 0xd4, 0x30,
	0x10, 0x55, 0x9a, 0x4d, 0x76, 0x77, 0x4a, 0xab, 0x62, 0xda, 0xc5, 0xaa, 0x00, 0x45, 0x39, 0x40,
	0x04, 0xd2, 0x4a, 0x2d, 0x5f, 0xe2, 0x84, 0xda, 0x82, 0x04, 0x12, 0x42, 0x60, 0x21, 0x2e, 0x9c,
	0xdc, 0x64, 0xe8, 0x46, 0xcd, 0xda, 0xc1, 0x71, 0x56, 0xec, 0x99, 0xff, 0xc2, 0x7f, 0xe2, 0xa7,
	0x70, 0x43, 0x76, 0x9c, 0x6e, 0xf6, 0x43, 0x2a, 0xbd, 0xf9, 0x8d, 0xdf, 0x8c, 0x3c, 0xcf, 0xf3,
	0x06, 0x6e, 0xa5, 0x13, 0x2e, 0x2e, 0x70, 0x5c, 0x2a, 0xa9, 0x25, 0x09, 0xf4, 0xbc, 0xc4, 0xea,
	0x70, 0xa7, 0x48, 0xe5, 0x74, 0x2a, 0x45, 0x13, 0x8d, 0xff, 0x6c, 0xc1, 0xe8, 0xa4, 0xd6, 0x52,
	0xc8, 0xe9, 0xfc, 0x93, 0x92, 0xa5, 0xac, 0x78, 0x71, 0x66, 0xd3, 0xc8, 0x73, 0x80, 0x52, 0xc9,
	0xb2, 0x41, 0xd4, 0x8b, 0xbc, 0x64, 0xfb, 0xf8, 0x60, 0x6c, 0xab, 0x8c, 0x97, 0xa9, 0xac, 0x43,
	0x24, 0x4f, 0xa0, 0x9f, 0xd6, 0x8a, 0xd5, 0x05, 0xd2, 0x2d, 0x9b, 0x73, 0xdb, 0xe5, 0x98, 0xd0,
	0x99, 0x14, 0xdf, 0xf3, 0x0b, 0xd6, 0x32, 0x48, 0x02, 0xc1, 0xb9, 0xe4, 0x2a, 0xa3, 0xbe, 0xa5,
	0x12, 0x47, 0x3d, 0x49, 0x75, 0x3e, 0xc3, 0x53, 0x73, 0xc3, 0x1a, 0x02, 0x39, 0x02, 0x98, 0x49,
	0x8d, 0x0c, 0xab, 0xba, 0xd0, 0xb4, 0xb7, 0x54, 0xf9, 0xeb, 0xd5, 0x05, 0xeb, 0x90, 0xc8, 0x08,
	0xc2, 0x4a, 0x73, 0x5d, 0x57, 0x34, 0x88, 0xbc, 0x24, 0x60, 0x0e, 0x11, 0x0a, 0x7d, 0x9e, 0x65,
	0x0a, 0xab, 0x8a, 0x86, 0x91, 0x97, 0x0c, 0x59, 0x0b, 0x4d, 0xc6, 0x04, 0xf3, 0x8b, 0x89, 0xa6,
	0xfd, 0xc8, 0x4b, 0x7c, 0xe6, 0x10, 0xd9, 0x87, 0x20, 0x17, 0x19, 0xfe, 0xa4, 0x03, 0x5b, 0xa8,
	0x01, 0xe4, 0x41, 0x23, 0x90, 0xd1, 0xe1, 0xfd, 0x1b, 0x3a, 0xb4, 0xa5, 0x3a, 0x91, 0xf8, 0xaf,
	0x07, 0xbb, 0x2b, 0x9a, 0x12, 0xe8, 0xcd, 0x91, 0x2b, 0xab, 0x66, 0xc0, 0xec, 0xd9, 0x14, 0x9f,
	0x4a, 0xa1, 0x27, 0x56, 0xae, 0x80, 0x35, 0x80, 0xec, 0x81, 0x9f, 0xf1, 0xb9, 0xd5, 0x25, 0x60,
	0xe6, 0x48, 0x1e, 0x41, 0xbf, 0xf9, 0xd0, 0x8a, 0xf6, 0x22, 0x3f, 0xd9, 0x3e, 0xde, 0x71, 0xed,
	0xbb, 0x4f, 0x68, 0x6f, 0xc9, 0x63, 0xd8, 0xab, 0x34, 0x57, 0xfa, 0xb4, 0x90, 0xe9, 0xe5, 0xbb,
	0xa6, 0x9f, 0xc0, 0xf6, 0xb3, 0x16, 0x27, 0x0f, 0x61, 0x17, 0x45, 0xd6, 0x65, 0x86, 0x96, 0xb9,
	0x12, 0x25, 0x63, 0x20, 0x0a, 0x79, 0xf1, 0x76, 0x99, 0xdb, 0xa8, 0xb4, 0xe1, 0x26, 0x7e, 0x06,
	0xa1, 0x6b, 0x79, 0x04, 0x61, 0xca, 0x45, 0x8a, 0x85, 0x6d, 0x7a, 0xc0, 0x1c, 0x32, 0x52, 0x18,
	0xd9, 0x6d, 0xd7, 0x43, 0x66, 0xcf, 0xf1, 0x0b, 0xd8, 0x67, 0x38, 0x93, 0x97, 0xb8, 0x22, 0xdb,
	0xb2, 0xd2, 0xde, 0x9a, 0xd2, 0x1f, 0x81, 0x98, 0x19, 0xb8, 0x59, 0x96, 0x9d, 0x83, 0xb2, 0x54,
	0x72, 0xd6, 0x4c, 0xea, 0x80, 0xb5, 0x30, 0x7e, 0x05, 0x77, 0xbf, 0xa0, 0x9a, 0xe6, 0x82, 0xdf,
	0xb4, 0x68, 0xfc, 0xcb, 0x83, 0x03, 0x86, 0x29, 0xe6, 0xa5, 0x5e, 0xc9, 0x3c, 0x82, 0x5e, 0xa9,
	0x70, 0xe6, 0x9c, 0x74, 0xbf, 0x1d, 0xf5, 0x8d, 0xe6, 0x63, 0x96, 0x4a, 0x5e, 0x5a, 0x2f, 0x29,
	0x14, 0xda, 0x79, 0xe9, 0x9a, 0xac, 0x96, 0x1d, 0x4f, 0xe0, 0xce, 0x07, 0x99, 0xf2, 0x62, 0xcd,
	0xd2, 0xa1, 0x79, 0xea, 0x69, 0xf6, 0x7f, 0x8f, 0x70, 0x64, 0x72, 0x08, 0x03, 0xb3, 0x34, 0x50,
	0xe8, 0x8a, 0x6e, 0x45, 0x7e, 0x32, 0x64, 0x57, 0x38, 0xfe, 0xed, 0xc1, 0x88, 0xe1, 0x8f, 0xcf,
	0x35, 0xaa, 0xd5, 0x05, 0xb2, 0xf0, 0x9f, 0xb7, 0xe4, 0xbf, 0x0d, 0x3f, 0x6f, 0x4c, 0x90, 0xca,
	0x5a, 0x68, 0x37, 0xf0, 0x0d, 0x20, 0xf7, 0x60, 0x98, 0xe5, 0x0a, 0x53, 0x9d, 0x4b, 0x61, 0x3d,
	0x1f, 0xb0, 0x45, 0xa0, 0xe3, 0xd6, 0x60, 0xb3, 0x5b, 0xc3, 0x8e, 0x5b, 0xe3, 0x6f, 0x40, 0x19,
	0x96, 0xc5, 0x7c, 0xd3, 0x4b, 0x5f, 0xc3, 0xf6, 0x62, 0x83, 0x99, 0xe7, 0xfa, 0xd7, 0x8b, 0xd3,
	0xcd, 0x38, 0x0f, 0xed, 0x36, 0x7d, 0xfa, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x72, 0xff, 0xf5, 0x9d,
	0x73, 0x05, 0x00, 0x00,
}
