package service

import (
	context "context"
)

const (
	CancelConnCmd = 1
	SendMsgCmd    = 2
)

type CmdContext struct {
	Ctx      *context.Context
	Cmd      int32
	Endpoint string
	ConnID   uint64
	Payload  []byte
}

type Service struct {
	CmdChannel chan *CmdContext
}

func (*Service) mustEmbedUnimplementedStateServer() {
	panic("unimplemented")
}

func (s *Service) CancelConn(ctx context.Context, sr *StateRequest) (*StateResponse, error) {
	c := context.TODO()
	s.CmdChannel <- &CmdContext{
		Ctx:      &c,
		Cmd:      CancelConnCmd,
		ConnID:   sr.ConnID,
		Endpoint: sr.GetEndpoint(),
	}
	return &StateResponse{
		Code: 0,
		Msg:  "success",
	}, nil
}

func (s *Service) SendMsg(ctx context.Context, sr *StateRequest) (*StateResponse, error) {
	c := context.TODO()
	s.CmdChannel <- &CmdContext{
		Ctx:      &c,
		Cmd:      SendMsgCmd,
		ConnID:   sr.ConnID,
		Endpoint: sr.GetEndpoint(),
		Payload:  sr.GetData(),
	}
	return &StateResponse{
		Code: 0,
		Msg:  "success",
	}, nil
}
