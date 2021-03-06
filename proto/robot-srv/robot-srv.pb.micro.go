// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/robot-srv/robot-srv.proto

package robot_srv

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for RobotSrv service

func NewRobotSrvEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for RobotSrv service

type RobotSrvService interface {
	SendMsg(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Test(ctx context.Context, in *TestRequest, opts ...client.CallOption) (*Response, error)
}

type robotSrvService struct {
	c    client.Client
	name string
}

func NewRobotSrvService(name string, c client.Client) RobotSrvService {
	return &robotSrvService{
		c:    c,
		name: name,
	}
}

func (c *robotSrvService) SendMsg(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "RobotSrv.SendMsg", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *robotSrvService) Test(ctx context.Context, in *TestRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "RobotSrv.Test", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RobotSrv service

type RobotSrvHandler interface {
	SendMsg(context.Context, *Request, *Response) error
	Test(context.Context, *TestRequest, *Response) error
}

func RegisterRobotSrvHandler(s server.Server, hdlr RobotSrvHandler, opts ...server.HandlerOption) error {
	type robotSrv interface {
		SendMsg(ctx context.Context, in *Request, out *Response) error
		Test(ctx context.Context, in *TestRequest, out *Response) error
	}
	type RobotSrv struct {
		robotSrv
	}
	h := &robotSrvHandler{hdlr}
	return s.Handle(s.NewHandler(&RobotSrv{h}, opts...))
}

type robotSrvHandler struct {
	RobotSrvHandler
}

func (h *robotSrvHandler) SendMsg(ctx context.Context, in *Request, out *Response) error {
	return h.RobotSrvHandler.SendMsg(ctx, in, out)
}

func (h *robotSrvHandler) Test(ctx context.Context, in *TestRequest, out *Response) error {
	return h.RobotSrvHandler.Test(ctx, in, out)
}
