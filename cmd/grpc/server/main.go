package main

import (
	"awesomeProject/accounts/models"
	"awesomeProject/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"sync"
)

func New() *server {
	return &server{
		accounts: make(map[string]*models.Account),
		guard:    &sync.RWMutex{},
	}
}

type Handler struct {
}

type server struct {
	proto.UnimplementedAccountServer
	accounts map[string]*models.Account
	guard    *sync.RWMutex
}

func (s *server) Get(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountReply, error) {
	s.guard.RLock()
	account, ok := s.accounts[req.GetName()]
	s.guard.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	return &proto.GetAccountReply{Name: account.Name, Amount: int32(account.Amount)}, nil
}

func (s *server) Create(ctx context.Context, req *proto.CreateAccountRequest) (*proto.Empty, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}

	s.guard.Lock()
	if _, ok := s.accounts[req.GetName()]; ok {
		s.guard.Unlock()

		return nil, status.Errorf(codes.AlreadyExists, "account already exists")
	}
	s.accounts[req.GetName()] = &models.Account{
		Name:   req.GetName(),
		Amount: int(req.GetAmount()),
	}
	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) ChangeAmount(ctx context.Context, req *proto.PatchAccountRequest) (*proto.Empty, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}

	s.guard.RLock()
	account, ok := s.accounts[req.GetName()]
	s.guard.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	s.guard.Lock()
	account.Amount = int(req.GetAmount())
	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) ChangeName(ctx context.Context, req *proto.ChangeAccountRequest) (*proto.Empty, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	if len(req.GetNewName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty new name")
	}

	s.guard.RLock()
	account, ok := s.accounts[req.GetName()]
	s.guard.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	s.guard.Lock()
	delete(s.accounts, account.Name)
	account.Name = req.GetNewName()
	s.accounts[req.GetNewName()] = account
	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.Empty, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}

	s.guard.RLock()
	account, ok := s.accounts[req.GetName()]
	s.guard.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	s.guard.Lock()
	delete(s.accounts, account.Name)
	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4567))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterAccountServer(s, New())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}

}
