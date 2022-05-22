// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: media.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Video struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkCids  []string `protobuf:"bytes,1,rep,name=chunkCids,proto3" json:"chunkCids,omitempty"`
	Duration   int64    `protobuf:"varint,2,opt,name=duration,proto3" json:"duration,omitempty"`
	FrameRate  int64    `protobuf:"varint,3,opt,name=frameRate,proto3" json:"frameRate,omitempty"`
	SampleRate int64    `protobuf:"varint,4,opt,name=sampleRate,proto3" json:"sampleRate,omitempty"`
}

func (x *Video) Reset() {
	*x = Video{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Video) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Video) ProtoMessage() {}

func (x *Video) ProtoReflect() protoreflect.Message {
	mi := &file_media_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Video.ProtoReflect.Descriptor instead.
func (*Video) Descriptor() ([]byte, []int) {
	return file_media_proto_rawDescGZIP(), []int{0}
}

func (x *Video) GetChunkCids() []string {
	if x != nil {
		return x.ChunkCids
	}
	return nil
}

func (x *Video) GetDuration() int64 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *Video) GetFrameRate() int64 {
	if x != nil {
		return x.FrameRate
	}
	return 0
}

func (x *Video) GetSampleRate() int64 {
	if x != nil {
		return x.SampleRate
	}
	return 0
}

type ChunkVideo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Frames  []*ImageRGBA   `protobuf:"bytes,1,rep,name=frames,proto3" json:"frames,omitempty"`
	Samples []*AudioSample `protobuf:"bytes,2,rep,name=samples,proto3" json:"samples,omitempty"`
}

func (x *ChunkVideo) Reset() {
	*x = ChunkVideo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChunkVideo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChunkVideo) ProtoMessage() {}

func (x *ChunkVideo) ProtoReflect() protoreflect.Message {
	mi := &file_media_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChunkVideo.ProtoReflect.Descriptor instead.
func (*ChunkVideo) Descriptor() ([]byte, []int) {
	return file_media_proto_rawDescGZIP(), []int{1}
}

func (x *ChunkVideo) GetFrames() []*ImageRGBA {
	if x != nil {
		return x.Frames
	}
	return nil
}

func (x *ChunkVideo) GetSamples() []*AudioSample {
	if x != nil {
		return x.Samples
	}
	return nil
}

type Audio struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkCids  []string `protobuf:"bytes,1,rep,name=chunkCids,proto3" json:"chunkCids,omitempty"`
	Duration   int64    `protobuf:"varint,2,opt,name=duration,proto3" json:"duration,omitempty"`
	SampleRate int64    `protobuf:"varint,3,opt,name=sampleRate,proto3" json:"sampleRate,omitempty"`
}

func (x *Audio) Reset() {
	*x = Audio{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Audio) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Audio) ProtoMessage() {}

func (x *Audio) ProtoReflect() protoreflect.Message {
	mi := &file_media_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Audio.ProtoReflect.Descriptor instead.
func (*Audio) Descriptor() ([]byte, []int) {
	return file_media_proto_rawDescGZIP(), []int{2}
}

func (x *Audio) GetChunkCids() []string {
	if x != nil {
		return x.ChunkCids
	}
	return nil
}

func (x *Audio) GetDuration() int64 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *Audio) GetSampleRate() int64 {
	if x != nil {
		return x.SampleRate
	}
	return 0
}

type ChunkAudio struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Samples []*AudioSample `protobuf:"bytes,1,rep,name=samples,proto3" json:"samples,omitempty"`
}

func (x *ChunkAudio) Reset() {
	*x = ChunkAudio{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChunkAudio) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChunkAudio) ProtoMessage() {}

func (x *ChunkAudio) ProtoReflect() protoreflect.Message {
	mi := &file_media_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChunkAudio.ProtoReflect.Descriptor instead.
func (*ChunkAudio) Descriptor() ([]byte, []int) {
	return file_media_proto_rawDescGZIP(), []int{3}
}

func (x *ChunkAudio) GetSamples() []*AudioSample {
	if x != nil {
		return x.Samples
	}
	return nil
}

type ImageRGBA struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pix    []uint32   `protobuf:"varint,1,rep,packed,name=pix,proto3" json:"pix,omitempty"`
	Stride int64      `protobuf:"varint,2,opt,name=stride,proto3" json:"stride,omitempty"`
	Rect   *Rectangle `protobuf:"bytes,3,opt,name=rect,proto3" json:"rect,omitempty"`
}

func (x *ImageRGBA) Reset() {
	*x = ImageRGBA{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageRGBA) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageRGBA) ProtoMessage() {}

func (x *ImageRGBA) ProtoReflect() protoreflect.Message {
	mi := &file_media_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageRGBA.ProtoReflect.Descriptor instead.
func (*ImageRGBA) Descriptor() ([]byte, []int) {
	return file_media_proto_rawDescGZIP(), []int{4}
}

func (x *ImageRGBA) GetPix() []uint32 {
	if x != nil {
		return x.Pix
	}
	return nil
}

func (x *ImageRGBA) GetStride() int64 {
	if x != nil {
		return x.Stride
	}
	return 0
}

func (x *ImageRGBA) GetRect() *Rectangle {
	if x != nil {
		return x.Rect
	}
	return nil
}

type Rectangle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Min *Point `protobuf:"bytes,1,opt,name=min,proto3" json:"min,omitempty"`
	Max *Point `protobuf:"bytes,2,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *Rectangle) Reset() {
	*x = Rectangle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rectangle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rectangle) ProtoMessage() {}

func (x *Rectangle) ProtoReflect() protoreflect.Message {
	mi := &file_media_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rectangle.ProtoReflect.Descriptor instead.
func (*Rectangle) Descriptor() ([]byte, []int) {
	return file_media_proto_rawDescGZIP(), []int{5}
}

func (x *Rectangle) GetMin() *Point {
	if x != nil {
		return x.Min
	}
	return nil
}

func (x *Rectangle) GetMax() *Point {
	if x != nil {
		return x.Max
	}
	return nil
}

type Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X int64 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	Y int64 `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *Point) Reset() {
	*x = Point{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Point) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Point) ProtoMessage() {}

func (x *Point) ProtoReflect() protoreflect.Message {
	mi := &file_media_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Point.ProtoReflect.Descriptor instead.
func (*Point) Descriptor() ([]byte, []int) {
	return file_media_proto_rawDescGZIP(), []int{6}
}

func (x *Point) GetX() int64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Point) GetY() int64 {
	if x != nil {
		return x.Y
	}
	return 0
}

type AudioSample struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Left  float64 `protobuf:"fixed64,1,opt,name=left,proto3" json:"left,omitempty"`
	Right float64 `protobuf:"fixed64,2,opt,name=right,proto3" json:"right,omitempty"`
}

func (x *AudioSample) Reset() {
	*x = AudioSample{}
	if protoimpl.UnsafeEnabled {
		mi := &file_media_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AudioSample) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AudioSample) ProtoMessage() {}

func (x *AudioSample) ProtoReflect() protoreflect.Message {
	mi := &file_media_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AudioSample.ProtoReflect.Descriptor instead.
func (*AudioSample) Descriptor() ([]byte, []int) {
	return file_media_proto_rawDescGZIP(), []int{7}
}

func (x *AudioSample) GetLeft() float64 {
	if x != nil {
		return x.Left
	}
	return 0
}

func (x *AudioSample) GetRight() float64 {
	if x != nil {
		return x.Right
	}
	return 0
}

var File_media_proto protoreflect.FileDescriptor

var file_media_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x22, 0x7f, 0x0a, 0x05, 0x56, 0x69, 0x64, 0x65, 0x6f,
	0x12, 0x1c, 0x0a, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x43, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x43, 0x69, 0x64, 0x73, 0x12, 0x1a,
	0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x72,
	0x61, 0x6d, 0x65, 0x52, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x66,
	0x72, 0x61, 0x6d, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x52, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x73, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x52, 0x61, 0x74, 0x65, 0x22, 0x6a, 0x0a, 0x0a, 0x43, 0x68, 0x75, 0x6e,
	0x6b, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x12, 0x2b, 0x0a, 0x06, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70,
	0x62, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x47, 0x42, 0x41, 0x52, 0x06, 0x66, 0x72, 0x61,
	0x6d, 0x65, 0x73, 0x12, 0x2f, 0x0a, 0x07, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e,
	0x41, 0x75, 0x64, 0x69, 0x6f, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x07, 0x73, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x73, 0x22, 0x61, 0x0a, 0x05, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x12, 0x1c, 0x0a,
	0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x43, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x43, 0x69, 0x64, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x52, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x73, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x52, 0x61, 0x74, 0x65, 0x22, 0x3d, 0x0a, 0x0a, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x41, 0x75, 0x64, 0x69, 0x6f, 0x12, 0x2f, 0x0a, 0x07, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70,
	0x62, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x07, 0x73,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x22, 0x5e, 0x0a, 0x09, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52,
	0x47, 0x42, 0x41, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x78, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d,
	0x52, 0x03, 0x70, 0x69, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x69, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x74, 0x72, 0x69, 0x64, 0x65, 0x12, 0x27, 0x0a,
	0x04, 0x72, 0x65, 0x63, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x63, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65,
	0x52, 0x04, 0x72, 0x65, 0x63, 0x74, 0x22, 0x51, 0x0a, 0x09, 0x52, 0x65, 0x63, 0x74, 0x61, 0x6e,
	0x67, 0x6c, 0x65, 0x12, 0x21, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x21, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x62, 0x2e, 0x50,
	0x6f, 0x69, 0x6e, 0x74, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x22, 0x23, 0x0a, 0x05, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x78,
	0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x79, 0x22, 0x37,
	0x0a, 0x0b, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x65, 0x66, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6c, 0x65, 0x66,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x05, 0x72, 0x69, 0x67, 0x68, 0x74, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_media_proto_rawDescOnce sync.Once
	file_media_proto_rawDescData = file_media_proto_rawDesc
)

func file_media_proto_rawDescGZIP() []byte {
	file_media_proto_rawDescOnce.Do(func() {
		file_media_proto_rawDescData = protoimpl.X.CompressGZIP(file_media_proto_rawDescData)
	})
	return file_media_proto_rawDescData
}

var file_media_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_media_proto_goTypes = []interface{}{
	(*Video)(nil),       // 0: store.pb.Video
	(*ChunkVideo)(nil),  // 1: store.pb.ChunkVideo
	(*Audio)(nil),       // 2: store.pb.Audio
	(*ChunkAudio)(nil),  // 3: store.pb.ChunkAudio
	(*ImageRGBA)(nil),   // 4: store.pb.ImageRGBA
	(*Rectangle)(nil),   // 5: store.pb.Rectangle
	(*Point)(nil),       // 6: store.pb.Point
	(*AudioSample)(nil), // 7: store.pb.AudioSample
}
var file_media_proto_depIdxs = []int32{
	4, // 0: store.pb.ChunkVideo.frames:type_name -> store.pb.ImageRGBA
	7, // 1: store.pb.ChunkVideo.samples:type_name -> store.pb.AudioSample
	7, // 2: store.pb.ChunkAudio.samples:type_name -> store.pb.AudioSample
	5, // 3: store.pb.ImageRGBA.rect:type_name -> store.pb.Rectangle
	6, // 4: store.pb.Rectangle.min:type_name -> store.pb.Point
	6, // 5: store.pb.Rectangle.max:type_name -> store.pb.Point
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_media_proto_init() }
func file_media_proto_init() {
	if File_media_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_media_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Video); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_media_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChunkVideo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_media_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Audio); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_media_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChunkAudio); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_media_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageRGBA); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_media_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rectangle); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_media_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Point); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_media_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AudioSample); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_media_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_media_proto_goTypes,
		DependencyIndexes: file_media_proto_depIdxs,
		MessageInfos:      file_media_proto_msgTypes,
	}.Build()
	File_media_proto = out.File
	file_media_proto_rawDesc = nil
	file_media_proto_goTypes = nil
	file_media_proto_depIdxs = nil
}
