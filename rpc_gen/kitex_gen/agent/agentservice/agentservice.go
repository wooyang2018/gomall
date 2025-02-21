// Code generated by Kitex v0.9.1. DO NOT EDIT.

package agentservice

import (
	"context"
	"errors"
	agent "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"QueryOrder": kitex.NewMethodInfo(
		queryOrderHandler,
		newQueryOrderArgs,
		newQueryOrderResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
}

var (
	agentServiceServiceInfo                = NewServiceInfo()
	agentServiceServiceInfoForClient       = NewServiceInfoForClient()
	agentServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return agentServiceServiceInfo
}

// for client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return agentServiceServiceInfoForStreamClient
}

// for stream client
func serviceInfoForClient() *kitex.ServiceInfo {
	return agentServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "AgentService"
	handlerType := (*agent.AgentService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "agent",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.9.1",
		Extra:           extra,
	}
	return svcInfo
}

func queryOrderHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(agent.QueryOrderReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(agent.AgentService).QueryOrder(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *QueryOrderArgs:
		success, err := handler.(agent.AgentService).QueryOrder(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*QueryOrderResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newQueryOrderArgs() interface{} {
	return &QueryOrderArgs{}
}

func newQueryOrderResult() interface{} {
	return &QueryOrderResult{}
}

type QueryOrderArgs struct {
	Req *agent.QueryOrderReq
}

func (p *QueryOrderArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(agent.QueryOrderReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *QueryOrderArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *QueryOrderArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *QueryOrderArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *QueryOrderArgs) Unmarshal(in []byte) error {
	msg := new(agent.QueryOrderReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var QueryOrderArgs_Req_DEFAULT *agent.QueryOrderReq

func (p *QueryOrderArgs) GetReq() *agent.QueryOrderReq {
	if !p.IsSetReq() {
		return QueryOrderArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *QueryOrderArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *QueryOrderArgs) GetFirstArgument() interface{} {
	return p.Req
}

type QueryOrderResult struct {
	Success *agent.QueryOrderResp
}

var QueryOrderResult_Success_DEFAULT *agent.QueryOrderResp

func (p *QueryOrderResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(agent.QueryOrderResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *QueryOrderResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *QueryOrderResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *QueryOrderResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *QueryOrderResult) Unmarshal(in []byte) error {
	msg := new(agent.QueryOrderResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *QueryOrderResult) GetSuccess() *agent.QueryOrderResp {
	if !p.IsSetSuccess() {
		return QueryOrderResult_Success_DEFAULT
	}
	return p.Success
}

func (p *QueryOrderResult) SetSuccess(x interface{}) {
	p.Success = x.(*agent.QueryOrderResp)
}

func (p *QueryOrderResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *QueryOrderResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) QueryOrder(ctx context.Context, Req *agent.QueryOrderReq) (r *agent.QueryOrderResp, err error) {
	var _args QueryOrderArgs
	_args.Req = Req
	var _result QueryOrderResult
	if err = p.c.Call(ctx, "QueryOrder", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
