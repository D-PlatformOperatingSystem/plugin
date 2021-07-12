// Code generated by protoc-gen-go. DO NOT EDIT.
// source: accountmanager.proto

package types

import (
	context "context"
	fmt "fmt"
	math "math"

	types "github.com/D-PlatformOperatingSystem/dpos/types"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type Accountmanager struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Accountmanager) Reset()         { *m = Accountmanager{} }
func (m *Accountmanager) String() string { return proto.CompactTextString(m) }
func (*Accountmanager) ProtoMessage()    {}
func (*Accountmanager) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{0}
}

func (m *Accountmanager) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Accountmanager.Unmarshal(m, b)
}
func (m *Accountmanager) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Accountmanager.Marshal(b, m, deterministic)
}
func (m *Accountmanager) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Accountmanager.Merge(m, src)
}
func (m *Accountmanager) XXX_Size() int {
	return xxx_messageInfo_Accountmanager.Size(m)
}
func (m *Accountmanager) XXX_DiscardUnknown() {
	xxx_messageInfo_Accountmanager.DiscardUnknown(m)
}

var xxx_messageInfo_Accountmanager proto.InternalMessageInfo

type AccountmanagerAction struct {
	// Types that are valid to be assigned to Value:
	//	*AccountmanagerAction_Register
	//	*AccountmanagerAction_ResetKey
	//	*AccountmanagerAction_Transfer
	//	*AccountmanagerAction_Supervise
	//	*AccountmanagerAction_Apply
	Value                isAccountmanagerAction_Value `protobuf_oneof:"value"`
	Ty                   int32                        `protobuf:"varint,6,opt,name=ty,proto3" json:"ty,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *AccountmanagerAction) Reset()         { *m = AccountmanagerAction{} }
func (m *AccountmanagerAction) String() string { return proto.CompactTextString(m) }
func (*AccountmanagerAction) ProtoMessage()    {}
func (*AccountmanagerAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{1}
}

func (m *AccountmanagerAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountmanagerAction.Unmarshal(m, b)
}
func (m *AccountmanagerAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountmanagerAction.Marshal(b, m, deterministic)
}
func (m *AccountmanagerAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountmanagerAction.Merge(m, src)
}
func (m *AccountmanagerAction) XXX_Size() int {
	return xxx_messageInfo_AccountmanagerAction.Size(m)
}
func (m *AccountmanagerAction) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountmanagerAction.DiscardUnknown(m)
}

var xxx_messageInfo_AccountmanagerAction proto.InternalMessageInfo

type isAccountmanagerAction_Value interface {
	isAccountmanagerAction_Value()
}

type AccountmanagerAction_Register struct {
	Register *Register `protobuf:"bytes,1,opt,name=register,proto3,oneof"`
}

type AccountmanagerAction_ResetKey struct {
	ResetKey *ResetKey `protobuf:"bytes,2,opt,name=resetKey,proto3,oneof"`
}

type AccountmanagerAction_Transfer struct {
	Transfer *Transfer `protobuf:"bytes,3,opt,name=transfer,proto3,oneof"`
}

type AccountmanagerAction_Supervise struct {
	Supervise *Supervise `protobuf:"bytes,4,opt,name=supervise,proto3,oneof"`
}

type AccountmanagerAction_Apply struct {
	Apply *Apply `protobuf:"bytes,5,opt,name=apply,proto3,oneof"`
}

func (*AccountmanagerAction_Register) isAccountmanagerAction_Value() {}

func (*AccountmanagerAction_ResetKey) isAccountmanagerAction_Value() {}

func (*AccountmanagerAction_Transfer) isAccountmanagerAction_Value() {}

func (*AccountmanagerAction_Supervise) isAccountmanagerAction_Value() {}

func (*AccountmanagerAction_Apply) isAccountmanagerAction_Value() {}

func (m *AccountmanagerAction) GetValue() isAccountmanagerAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *AccountmanagerAction) GetRegister() *Register {
	if x, ok := m.GetValue().(*AccountmanagerAction_Register); ok {
		return x.Register
	}
	return nil
}

func (m *AccountmanagerAction) GetResetKey() *ResetKey {
	if x, ok := m.GetValue().(*AccountmanagerAction_ResetKey); ok {
		return x.ResetKey
	}
	return nil
}

func (m *AccountmanagerAction) GetTransfer() *Transfer {
	if x, ok := m.GetValue().(*AccountmanagerAction_Transfer); ok {
		return x.Transfer
	}
	return nil
}

func (m *AccountmanagerAction) GetSupervise() *Supervise {
	if x, ok := m.GetValue().(*AccountmanagerAction_Supervise); ok {
		return x.Supervise
	}
	return nil
}

func (m *AccountmanagerAction) GetApply() *Apply {
	if x, ok := m.GetValue().(*AccountmanagerAction_Apply); ok {
		return x.Apply
	}
	return nil
}

func (m *AccountmanagerAction) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*AccountmanagerAction) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*AccountmanagerAction_Register)(nil),
		(*AccountmanagerAction_ResetKey)(nil),
		(*AccountmanagerAction_Transfer)(nil),
		(*AccountmanagerAction_Supervise)(nil),
		(*AccountmanagerAction_Apply)(nil),
	}
}

//
type Register struct {
	AccountID            string   `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Register) Reset()         { *m = Register{} }
func (m *Register) String() string { return proto.CompactTextString(m) }
func (*Register) ProtoMessage()    {}
func (*Register) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{2}
}

func (m *Register) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Register.Unmarshal(m, b)
}
func (m *Register) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Register.Marshal(b, m, deterministic)
}
func (m *Register) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Register.Merge(m, src)
}
func (m *Register) XXX_Size() int {
	return xxx_messageInfo_Register.Size(m)
}
func (m *Register) XXX_DiscardUnknown() {
	xxx_messageInfo_Register.DiscardUnknown(m)
}

var xxx_messageInfo_Register proto.InternalMessageInfo

func (m *Register) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

//
type ResetKey struct {
	AccountID            string   `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResetKey) Reset()         { *m = ResetKey{} }
func (m *ResetKey) String() string { return proto.CompactTextString(m) }
func (*ResetKey) ProtoMessage()    {}
func (*ResetKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{3}
}

func (m *ResetKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResetKey.Unmarshal(m, b)
}
func (m *ResetKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResetKey.Marshal(b, m, deterministic)
}
func (m *ResetKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResetKey.Merge(m, src)
}
func (m *ResetKey) XXX_Size() int {
	return xxx_messageInfo_ResetKey.Size(m)
}
func (m *ResetKey) XXX_DiscardUnknown() {
	xxx_messageInfo_ResetKey.DiscardUnknown(m)
}

var xxx_messageInfo_ResetKey proto.InternalMessageInfo

func (m *ResetKey) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *ResetKey) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

//
type Apply struct {
	AccountID string `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	//  ， 1         , 2       ，
	Op                   int32    `protobuf:"varint,2,opt,name=op,proto3" json:"op,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Apply) Reset()         { *m = Apply{} }
func (m *Apply) String() string { return proto.CompactTextString(m) }
func (*Apply) ProtoMessage()    {}
func (*Apply) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{4}
}

func (m *Apply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Apply.Unmarshal(m, b)
}
func (m *Apply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Apply.Marshal(b, m, deterministic)
}
func (m *Apply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Apply.Merge(m, src)
}
func (m *Apply) XXX_Size() int {
	return xxx_messageInfo_Apply.Size(m)
}
func (m *Apply) XXX_DiscardUnknown() {
	xxx_messageInfo_Apply.DiscardUnknown(m)
}

var xxx_messageInfo_Apply proto.InternalMessageInfo

func (m *Apply) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *Apply) GetOp() int32 {
	if m != nil {
		return m.Op
	}
	return 0
}

//
type Transfer struct {
	//
	Asset *types.Asset `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset,omitempty"`
	// from
	FromAccountID string `protobuf:"bytes,2,opt,name=fromAccountID,proto3" json:"fromAccountID,omitempty"`
	// to
	ToAccountID          string   `protobuf:"bytes,3,opt,name=toAccountID,proto3" json:"toAccountID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transfer) Reset()         { *m = Transfer{} }
func (m *Transfer) String() string { return proto.CompactTextString(m) }
func (*Transfer) ProtoMessage()    {}
func (*Transfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{5}
}

func (m *Transfer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transfer.Unmarshal(m, b)
}
func (m *Transfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transfer.Marshal(b, m, deterministic)
}
func (m *Transfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transfer.Merge(m, src)
}
func (m *Transfer) XXX_Size() int {
	return xxx_messageInfo_Transfer.Size(m)
}
func (m *Transfer) XXX_DiscardUnknown() {
	xxx_messageInfo_Transfer.DiscardUnknown(m)
}

var xxx_messageInfo_Transfer proto.InternalMessageInfo

func (m *Transfer) GetAsset() *types.Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

func (m *Transfer) GetFromAccountID() string {
	if m != nil {
		return m.FromAccountID
	}
	return ""
}

func (m *Transfer) GetToAccountID() string {
	if m != nil {
		return m.ToAccountID
	}
	return ""
}

//
type Supervise struct {
	//
	AccountIDs []string `protobuf:"bytes,1,rep,name=accountIDs,proto3" json:"accountIDs,omitempty"`
	//  ， 1   ，2   ，3     ,4
	Op int32 `protobuf:"varint,2,opt,name=op,proto3" json:"op,omitempty"`
	//0  ,             ，
	Level                int32    `protobuf:"varint,3,opt,name=level,proto3" json:"level,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Supervise) Reset()         { *m = Supervise{} }
func (m *Supervise) String() string { return proto.CompactTextString(m) }
func (*Supervise) ProtoMessage()    {}
func (*Supervise) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{6}
}

func (m *Supervise) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Supervise.Unmarshal(m, b)
}
func (m *Supervise) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Supervise.Marshal(b, m, deterministic)
}
func (m *Supervise) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Supervise.Merge(m, src)
}
func (m *Supervise) XXX_Size() int {
	return xxx_messageInfo_Supervise.Size(m)
}
func (m *Supervise) XXX_DiscardUnknown() {
	xxx_messageInfo_Supervise.DiscardUnknown(m)
}

var xxx_messageInfo_Supervise proto.InternalMessageInfo

func (m *Supervise) GetAccountIDs() []string {
	if m != nil {
		return m.AccountIDs
	}
	return nil
}

func (m *Supervise) GetOp() int32 {
	if m != nil {
		return m.Op
	}
	return 0
}

func (m *Supervise) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

type Account struct {
	//
	AccountID string `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	//
	Addr string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	//
	PrevAddr string `protobuf:"bytes,3,opt,name=prevAddr,proto3" json:"prevAddr,omitempty"`
	//     0   ， 1    , 2     3,
	Status int32 `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
	//     0  ,             ，
	Level int32 `protobuf:"varint,5,opt,name=level,proto3" json:"level,omitempty"`
	//
	CreateTime int64 `protobuf:"varint,6,opt,name=createTime,proto3" json:"createTime,omitempty"`
	//
	ExpireTime int64 `protobuf:"varint,7,opt,name=expireTime,proto3" json:"expireTime,omitempty"`
	//
	LockTime int64 `protobuf:"varint,8,opt,name=lockTime,proto3" json:"lockTime,omitempty"`
	//
	Index                int64    `protobuf:"varint,9,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{7}
}

func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *Account) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *Account) GetPrevAddr() string {
	if m != nil {
		return m.PrevAddr
	}
	return ""
}

func (m *Account) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *Account) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *Account) GetCreateTime() int64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *Account) GetExpireTime() int64 {
	if m != nil {
		return m.ExpireTime
	}
	return 0
}

func (m *Account) GetLockTime() int64 {
	if m != nil {
		return m.LockTime
	}
	return 0
}

func (m *Account) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

type AccountReceipt struct {
	Account              *Account `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountReceipt) Reset()         { *m = AccountReceipt{} }
func (m *AccountReceipt) String() string { return proto.CompactTextString(m) }
func (*AccountReceipt) ProtoMessage()    {}
func (*AccountReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{8}
}

func (m *AccountReceipt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountReceipt.Unmarshal(m, b)
}
func (m *AccountReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountReceipt.Marshal(b, m, deterministic)
}
func (m *AccountReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountReceipt.Merge(m, src)
}
func (m *AccountReceipt) XXX_Size() int {
	return xxx_messageInfo_AccountReceipt.Size(m)
}
func (m *AccountReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_AccountReceipt proto.InternalMessageInfo

func (m *AccountReceipt) GetAccount() *Account {
	if m != nil {
		return m.Account
	}
	return nil
}

type ReplyAccountList struct {
	Accounts             []*Account `protobuf:"bytes,1,rep,name=accounts,proto3" json:"accounts,omitempty"`
	PrimaryKey           string     `protobuf:"bytes,2,opt,name=primaryKey,proto3" json:"primaryKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ReplyAccountList) Reset()         { *m = ReplyAccountList{} }
func (m *ReplyAccountList) String() string { return proto.CompactTextString(m) }
func (*ReplyAccountList) ProtoMessage()    {}
func (*ReplyAccountList) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{9}
}

func (m *ReplyAccountList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReplyAccountList.Unmarshal(m, b)
}
func (m *ReplyAccountList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReplyAccountList.Marshal(b, m, deterministic)
}
func (m *ReplyAccountList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReplyAccountList.Merge(m, src)
}
func (m *ReplyAccountList) XXX_Size() int {
	return xxx_messageInfo_ReplyAccountList.Size(m)
}
func (m *ReplyAccountList) XXX_DiscardUnknown() {
	xxx_messageInfo_ReplyAccountList.DiscardUnknown(m)
}

var xxx_messageInfo_ReplyAccountList proto.InternalMessageInfo

func (m *ReplyAccountList) GetAccounts() []*Account {
	if m != nil {
		return m.Accounts
	}
	return nil
}

func (m *ReplyAccountList) GetPrimaryKey() string {
	if m != nil {
		return m.PrimaryKey
	}
	return ""
}

type TransferReceipt struct {
	FromAccount          *Account `protobuf:"bytes,1,opt,name=FromAccount,proto3" json:"FromAccount,omitempty"`
	ToAccount            *Account `protobuf:"bytes,2,opt,name=ToAccount,proto3" json:"ToAccount,omitempty"`
	Index                int64    `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransferReceipt) Reset()         { *m = TransferReceipt{} }
func (m *TransferReceipt) String() string { return proto.CompactTextString(m) }
func (*TransferReceipt) ProtoMessage()    {}
func (*TransferReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{10}
}

func (m *TransferReceipt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransferReceipt.Unmarshal(m, b)
}
func (m *TransferReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransferReceipt.Marshal(b, m, deterministic)
}
func (m *TransferReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferReceipt.Merge(m, src)
}
func (m *TransferReceipt) XXX_Size() int {
	return xxx_messageInfo_TransferReceipt.Size(m)
}
func (m *TransferReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_TransferReceipt proto.InternalMessageInfo

func (m *TransferReceipt) GetFromAccount() *Account {
	if m != nil {
		return m.FromAccount
	}
	return nil
}

func (m *TransferReceipt) GetToAccount() *Account {
	if m != nil {
		return m.ToAccount
	}
	return nil
}

func (m *TransferReceipt) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

//
type SuperviseReceipt struct {
	Accounts             []*Account `protobuf:"bytes,1,rep,name=accounts,proto3" json:"accounts,omitempty"`
	Op                   int32      `protobuf:"varint,2,opt,name=op,proto3" json:"op,omitempty"`
	Index                int64      `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *SuperviseReceipt) Reset()         { *m = SuperviseReceipt{} }
func (m *SuperviseReceipt) String() string { return proto.CompactTextString(m) }
func (*SuperviseReceipt) ProtoMessage()    {}
func (*SuperviseReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{11}
}

func (m *SuperviseReceipt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SuperviseReceipt.Unmarshal(m, b)
}
func (m *SuperviseReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SuperviseReceipt.Marshal(b, m, deterministic)
}
func (m *SuperviseReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SuperviseReceipt.Merge(m, src)
}
func (m *SuperviseReceipt) XXX_Size() int {
	return xxx_messageInfo_SuperviseReceipt.Size(m)
}
func (m *SuperviseReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_SuperviseReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_SuperviseReceipt proto.InternalMessageInfo

func (m *SuperviseReceipt) GetAccounts() []*Account {
	if m != nil {
		return m.Accounts
	}
	return nil
}

func (m *SuperviseReceipt) GetOp() int32 {
	if m != nil {
		return m.Op
	}
	return 0
}

func (m *SuperviseReceipt) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

type QueryExpiredAccounts struct {
	PrimaryKey string `protobuf:"bytes,1,opt,name=primaryKey,proto3" json:"primaryKey,omitempty"`
	//           ，
	ExpiredTime int64 `protobuf:"varint,2,opt,name=expiredTime,proto3" json:"expiredTime,omitempty"`
	//         ，    10
	// 0  ，1  ，
	Direction            int32    `protobuf:"varint,3,opt,name=direction,proto3" json:"direction,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryExpiredAccounts) Reset()         { *m = QueryExpiredAccounts{} }
func (m *QueryExpiredAccounts) String() string { return proto.CompactTextString(m) }
func (*QueryExpiredAccounts) ProtoMessage()    {}
func (*QueryExpiredAccounts) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{12}
}

func (m *QueryExpiredAccounts) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryExpiredAccounts.Unmarshal(m, b)
}
func (m *QueryExpiredAccounts) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryExpiredAccounts.Marshal(b, m, deterministic)
}
func (m *QueryExpiredAccounts) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryExpiredAccounts.Merge(m, src)
}
func (m *QueryExpiredAccounts) XXX_Size() int {
	return xxx_messageInfo_QueryExpiredAccounts.Size(m)
}
func (m *QueryExpiredAccounts) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryExpiredAccounts.DiscardUnknown(m)
}

var xxx_messageInfo_QueryExpiredAccounts proto.InternalMessageInfo

func (m *QueryExpiredAccounts) GetPrimaryKey() string {
	if m != nil {
		return m.PrimaryKey
	}
	return ""
}

func (m *QueryExpiredAccounts) GetExpiredTime() int64 {
	if m != nil {
		return m.ExpiredTime
	}
	return 0
}

func (m *QueryExpiredAccounts) GetDirection() int32 {
	if m != nil {
		return m.Direction
	}
	return 0
}

type QueryAccountsByStatus struct {
	//     1   ， 2    , 3
	Status int32 `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	//
	PrimaryKey string `protobuf:"bytes,3,opt,name=primaryKey,proto3" json:"primaryKey,omitempty"`
	// 0  ，1  ，
	Direction            int32    `protobuf:"varint,5,opt,name=direction,proto3" json:"direction,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryAccountsByStatus) Reset()         { *m = QueryAccountsByStatus{} }
func (m *QueryAccountsByStatus) String() string { return proto.CompactTextString(m) }
func (*QueryAccountsByStatus) ProtoMessage()    {}
func (*QueryAccountsByStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{13}
}

func (m *QueryAccountsByStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryAccountsByStatus.Unmarshal(m, b)
}
func (m *QueryAccountsByStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryAccountsByStatus.Marshal(b, m, deterministic)
}
func (m *QueryAccountsByStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAccountsByStatus.Merge(m, src)
}
func (m *QueryAccountsByStatus) XXX_Size() int {
	return xxx_messageInfo_QueryAccountsByStatus.Size(m)
}
func (m *QueryAccountsByStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAccountsByStatus.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAccountsByStatus proto.InternalMessageInfo

func (m *QueryAccountsByStatus) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *QueryAccountsByStatus) GetPrimaryKey() string {
	if m != nil {
		return m.PrimaryKey
	}
	return ""
}

func (m *QueryAccountsByStatus) GetDirection() int32 {
	if m != nil {
		return m.Direction
	}
	return 0
}

type QueryAccountByID struct {
	AccountID            string   `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryAccountByID) Reset()         { *m = QueryAccountByID{} }
func (m *QueryAccountByID) String() string { return proto.CompactTextString(m) }
func (*QueryAccountByID) ProtoMessage()    {}
func (*QueryAccountByID) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{14}
}

func (m *QueryAccountByID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryAccountByID.Unmarshal(m, b)
}
func (m *QueryAccountByID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryAccountByID.Marshal(b, m, deterministic)
}
func (m *QueryAccountByID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAccountByID.Merge(m, src)
}
func (m *QueryAccountByID) XXX_Size() int {
	return xxx_messageInfo_QueryAccountByID.Size(m)
}
func (m *QueryAccountByID) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAccountByID.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAccountByID proto.InternalMessageInfo

func (m *QueryAccountByID) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

type QueryAccountByAddr struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryAccountByAddr) Reset()         { *m = QueryAccountByAddr{} }
func (m *QueryAccountByAddr) String() string { return proto.CompactTextString(m) }
func (*QueryAccountByAddr) ProtoMessage()    {}
func (*QueryAccountByAddr) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{15}
}

func (m *QueryAccountByAddr) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryAccountByAddr.Unmarshal(m, b)
}
func (m *QueryAccountByAddr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryAccountByAddr.Marshal(b, m, deterministic)
}
func (m *QueryAccountByAddr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAccountByAddr.Merge(m, src)
}
func (m *QueryAccountByAddr) XXX_Size() int {
	return xxx_messageInfo_QueryAccountByAddr.Size(m)
}
func (m *QueryAccountByAddr) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAccountByAddr.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAccountByAddr proto.InternalMessageInfo

func (m *QueryAccountByAddr) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type QueryBalanceByID struct {
	AccountID            string       `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	Asset                *types.Asset `protobuf:"bytes,2,opt,name=asset,proto3" json:"asset,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *QueryBalanceByID) Reset()         { *m = QueryBalanceByID{} }
func (m *QueryBalanceByID) String() string { return proto.CompactTextString(m) }
func (*QueryBalanceByID) ProtoMessage()    {}
func (*QueryBalanceByID) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{16}
}

func (m *QueryBalanceByID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryBalanceByID.Unmarshal(m, b)
}
func (m *QueryBalanceByID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryBalanceByID.Marshal(b, m, deterministic)
}
func (m *QueryBalanceByID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBalanceByID.Merge(m, src)
}
func (m *QueryBalanceByID) XXX_Size() int {
	return xxx_messageInfo_QueryBalanceByID.Size(m)
}
func (m *QueryBalanceByID) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBalanceByID.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBalanceByID proto.InternalMessageInfo

func (m *QueryBalanceByID) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *QueryBalanceByID) GetAsset() *types.Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

type Balance struct {
	Balance              int64    `protobuf:"varint,1,opt,name=balance,proto3" json:"balance,omitempty"`
	Frozen               int64    `protobuf:"varint,2,opt,name=frozen,proto3" json:"frozen,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Balance) Reset()         { *m = Balance{} }
func (m *Balance) String() string { return proto.CompactTextString(m) }
func (*Balance) ProtoMessage()    {}
func (*Balance) Descriptor() ([]byte, []int) {
	return fileDescriptor_aba6db06a1aaf83a, []int{17}
}

func (m *Balance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Balance.Unmarshal(m, b)
}
func (m *Balance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Balance.Marshal(b, m, deterministic)
}
func (m *Balance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Balance.Merge(m, src)
}
func (m *Balance) XXX_Size() int {
	return xxx_messageInfo_Balance.Size(m)
}
func (m *Balance) XXX_DiscardUnknown() {
	xxx_messageInfo_Balance.DiscardUnknown(m)
}

var xxx_messageInfo_Balance proto.InternalMessageInfo

func (m *Balance) GetBalance() int64 {
	if m != nil {
		return m.Balance
	}
	return 0
}

func (m *Balance) GetFrozen() int64 {
	if m != nil {
		return m.Frozen
	}
	return 0
}

func init() {
	proto.RegisterType((*Accountmanager)(nil), "types.Accountmanager")
	proto.RegisterType((*AccountmanagerAction)(nil), "types.AccountmanagerAction")
	proto.RegisterType((*Register)(nil), "types.Register")
	proto.RegisterType((*ResetKey)(nil), "types.ResetKey")
	proto.RegisterType((*Apply)(nil), "types.Apply")
	proto.RegisterType((*Transfer)(nil), "types.Transfer")
	proto.RegisterType((*Supervise)(nil), "types.Supervise")
	proto.RegisterType((*Account)(nil), "types.account")
	proto.RegisterType((*AccountReceipt)(nil), "types.AccountReceipt")
	proto.RegisterType((*ReplyAccountList)(nil), "types.ReplyAccountList")
	proto.RegisterType((*TransferReceipt)(nil), "types.TransferReceipt")
	proto.RegisterType((*SuperviseReceipt)(nil), "types.SuperviseReceipt")
	proto.RegisterType((*QueryExpiredAccounts)(nil), "types.QueryExpiredAccounts")
	proto.RegisterType((*QueryAccountsByStatus)(nil), "types.QueryAccountsByStatus")
	proto.RegisterType((*QueryAccountByID)(nil), "types.QueryAccountByID")
	proto.RegisterType((*QueryAccountByAddr)(nil), "types.QueryAccountByAddr")
	proto.RegisterType((*QueryBalanceByID)(nil), "types.QueryBalanceByID")
	proto.RegisterType((*Balance)(nil), "types.balance")
}

func init() {
	proto.RegisterFile("accountmanager.proto", fileDescriptor_aba6db06a1aaf83a)
}

var fileDescriptor_aba6db06a1aaf83a = []byte{
	// 699 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x5d, 0x92, 0xa5, 0x6d, 0xee, 0xa0, 0x2b, 0x56, 0x99, 0xa2, 0x09, 0xa1, 0xca, 0xda, 0x43,
	0x85, 0x60, 0x9a, 0x86, 0x78, 0x01, 0x5e, 0x5a, 0x0d, 0xb4, 0x09, 0x5e, 0xe6, 0xf5, 0x19, 0xc9,
	0x4b, 0xee, 0xa6, 0x88, 0x34, 0x89, 0x1c, 0xb7, 0x5a, 0xf8, 0x03, 0xfc, 0x01, 0xfe, 0x2b, 0xaf,
	0x28, 0xb6, 0xf3, 0xd5, 0x6d, 0x54, 0x7b, 0x8b, 0xcf, 0x3d, 0xbe, 0xf7, 0xf8, 0xde, 0x63, 0x07,
	0xc6, 0x3c, 0x08, 0xd2, 0x55, 0x22, 0x97, 0x3c, 0xe1, 0xb7, 0x28, 0x8e, 0x33, 0x91, 0xca, 0x94,
	0xb8, 0xb2, 0xc8, 0x30, 0x3f, 0x7c, 0x21, 0x05, 0x4f, 0x72, 0x1e, 0xc8, 0x28, 0x4d, 0x74, 0x84,
	0x8e, 0x60, 0x38, 0xeb, 0xec, 0xa0, 0x7f, 0x6c, 0x18, 0x77, 0xa1, 0x99, 0xda, 0x40, 0xde, 0xc1,
	0x40, 0xe0, 0x6d, 0x94, 0x4b, 0x14, 0xbe, 0x35, 0xb1, 0xa6, 0x7b, 0xa7, 0xfb, 0xc7, 0x2a, 0xef,
	0x31, 0x33, 0xf0, 0xf9, 0x0e, 0xab, 0x29, 0x9a, 0x9e, 0xa3, 0xfc, 0x86, 0x85, 0x6f, 0x6f, 0xd0,
	0x35, 0xac, 0xe9, 0xfa, 0xbb, 0xa4, 0x2b, 0x75, 0x37, 0x28, 0x7c, 0xa7, 0x43, 0x5f, 0x18, 0xb8,
	0xa4, 0x57, 0x14, 0x72, 0x02, 0x5e, 0xbe, 0xca, 0x50, 0xac, 0xa3, 0x1c, 0xfd, 0x5d, 0xc5, 0x1f,
	0x19, 0xfe, 0x55, 0x85, 0x9f, 0xef, 0xb0, 0x86, 0x44, 0x8e, 0xc0, 0xe5, 0x59, 0x16, 0x17, 0xbe,
	0xab, 0xd8, 0xcf, 0x0c, 0x7b, 0x56, 0x62, 0xe7, 0x3b, 0x4c, 0x07, 0xc9, 0x10, 0x6c, 0x59, 0xf8,
	0xbd, 0x89, 0x35, 0x75, 0x99, 0x2d, 0x8b, 0x79, 0x1f, 0xdc, 0x35, 0x8f, 0x57, 0x48, 0xa7, 0x30,
	0xa8, 0x8e, 0x49, 0x5e, 0x81, 0x67, 0xda, 0x7c, 0x71, 0xa6, 0x5a, 0xe1, 0xb1, 0x06, 0xa0, 0x9f,
	0x4b, 0xa6, 0x39, 0xd5, 0x7f, 0x99, 0x84, 0xc0, 0x2e, 0x0f, 0x43, 0xa1, 0xda, 0xe3, 0x31, 0xf5,
	0x4d, 0x3f, 0x80, 0xab, 0x24, 0x6d, 0xd9, 0x3a, 0x04, 0x3b, 0xcd, 0xd4, 0x46, 0x97, 0xd9, 0x69,
	0x46, 0xd7, 0x30, 0xa8, 0xfa, 0x44, 0x28, 0xb8, 0x3c, 0xcf, 0x51, 0x9a, 0x29, 0xd5, 0x27, 0x2d,
	0x31, 0xa6, 0x43, 0xe4, 0x08, 0x9e, 0xdf, 0x88, 0x74, 0x39, 0xab, 0x2b, 0x68, 0x0d, 0x5d, 0x90,
	0x4c, 0x60, 0x4f, 0xa6, 0x0d, 0xc7, 0x51, 0x9c, 0x36, 0x44, 0x2f, 0xc1, 0xab, 0xfb, 0x4d, 0x5e,
	0x03, 0xd4, 0x0a, 0x73, 0xdf, 0x9a, 0x38, 0x53, 0x8f, 0xb5, 0x90, 0x4d, 0xd1, 0x64, 0x0c, 0x6e,
	0x8c, 0x6b, 0x8c, 0x55, 0x62, 0x97, 0xe9, 0x05, 0xfd, 0x6b, 0x41, 0xdf, 0x6c, 0x7a, 0x7a, 0xff,
	0xc8, 0x21, 0x0c, 0x32, 0x81, 0xeb, 0x59, 0x89, 0x6b, 0xbd, 0xf5, 0x9a, 0x1c, 0x40, 0x2f, 0x97,
	0x5c, 0xae, 0x72, 0xe5, 0x18, 0x97, 0x99, 0x55, 0xa3, 0xc3, 0x6d, 0xe9, 0x28, 0x4f, 0x13, 0x08,
	0xe4, 0x12, 0x17, 0xd1, 0x12, 0x95, 0x25, 0x1c, 0xd6, 0x42, 0xca, 0x38, 0xde, 0x65, 0x91, 0xd0,
	0xf1, 0xbe, 0x8e, 0x37, 0x48, 0xa9, 0x24, 0x4e, 0x83, 0x9f, 0x2a, 0x3a, 0x50, 0xd1, 0x7a, 0x5d,
	0x56, 0x8c, 0x92, 0x10, 0xef, 0x7c, 0x4f, 0x05, 0xf4, 0x82, 0x7e, 0xac, 0x2f, 0x23, 0xc3, 0x00,
	0xa3, 0x4c, 0x92, 0x69, 0xdd, 0x0a, 0x33, 0xcc, 0xa1, 0x19, 0xa6, 0x41, 0x59, 0x15, 0xa6, 0x3f,
	0x60, 0xc4, 0x30, 0x8b, 0x0b, 0x93, 0xe0, 0x7b, 0x94, 0x4b, 0xf2, 0x06, 0x06, 0x26, 0xac, 0xa7,
	0x71, 0x7f, 0x7b, 0x1d, 0x2f, 0x4f, 0x93, 0x89, 0x68, 0xc9, 0x45, 0x51, 0x5d, 0x58, 0x8f, 0xb5,
	0x10, 0xfa, 0xdb, 0x82, 0xfd, 0xca, 0x61, 0x95, 0xba, 0x13, 0xd8, 0xfb, 0xda, 0xf8, 0xe5, 0x11,
	0x85, 0x6d, 0x0a, 0x79, 0x0b, 0xde, 0xa2, 0x72, 0x8f, 0x79, 0x15, 0x36, 0xf9, 0x0d, 0xa1, 0xe9,
	0x92, 0xd3, 0xee, 0x52, 0x08, 0xa3, 0xda, 0x72, 0x95, 0x92, 0xa7, 0x9c, 0xf4, 0x01, 0x17, 0x3e,
	0x50, 0x65, 0x0d, 0xe3, 0xcb, 0x15, 0x8a, 0xe2, 0x8b, 0x1a, 0x68, 0x38, 0x7b, 0xb8, 0x4f, 0xd6,
	0x66, 0x9f, 0xca, 0x2b, 0xa3, 0x3d, 0x10, 0xaa, 0xc1, 0xdb, 0x2a, 0x67, 0x1b, 0x2a, 0x3d, 0x1d,
	0x46, 0x02, 0xd5, 0xa3, 0x6a, 0x9c, 0xdf, 0x00, 0x74, 0x09, 0x2f, 0x55, 0xdd, 0xaa, 0xe0, 0xbc,
	0xb8, 0xd2, 0x26, 0x6d, 0xcc, 0x6b, 0x75, 0xcc, 0xdb, 0x15, 0xe4, 0xdc, 0x13, 0xd4, 0x29, 0xe7,
	0x6e, 0x96, 0x3b, 0x81, 0x51, 0xbb, 0xdc, 0xbc, 0xb8, 0x38, 0xdb, 0xf2, 0xbc, 0x4d, 0x81, 0x74,
	0x77, 0xa8, 0xab, 0x55, 0x5d, 0x45, 0xab, 0xf5, 0x94, 0x2d, 0x4c, 0xee, 0x39, 0x8f, 0x79, 0x12,
	0xe0, 0xf6, 0xdc, 0xcd, 0xcb, 0x65, 0x3f, 0xfa, 0x72, 0xd1, 0x4f, 0xd0, 0xbf, 0xd6, 0x09, 0x89,
	0x5f, 0x7f, 0xaa, 0x54, 0x0e, 0xab, 0x23, 0x07, 0xd0, 0xbb, 0x11, 0xe9, 0x2f, 0x4c, 0xcc, 0x00,
	0xcc, 0xea, 0x74, 0x04, 0xc3, 0xee, 0x0f, 0xf2, 0xba, 0xa7, 0xfe, 0x83, 0xef, 0xff, 0x05, 0x00,
	0x00, 0xff, 0xff, 0x84, 0xc3, 0xe9, 0x87, 0x39, 0x07, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AccountmanagerClient is the client API for Accountmanager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AccountmanagerClient interface {
}

type accountmanagerClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountmanagerClient(cc grpc.ClientConnInterface) AccountmanagerClient {
	return &accountmanagerClient{cc}
}

// AccountmanagerServer is the server API for Accountmanager service.
type AccountmanagerServer interface {
}

// UnimplementedAccountmanagerServer can be embedded to have forward compatible implementations.
type UnimplementedAccountmanagerServer struct {
}

func RegisterAccountmanagerServer(s *grpc.Server, srv AccountmanagerServer) {
	s.RegisterService(&_Accountmanager_serviceDesc, srv)
}

var _Accountmanager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "types.accountmanager",
	HandlerType: (*AccountmanagerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "accountmanager.proto",
}