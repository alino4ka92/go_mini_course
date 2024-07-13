package main

import (
	"awesomeProject/accounts/models"
	"awesomeProject/proto"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

func New(db *sql.DB, ctx context.Context) *server {
	return &server{
		db:  db,
		ctx: ctx,
	}
}

type Handler struct {
}

type server struct {
	proto.UnimplementedAccountServer
	db  *sql.DB
	ctx context.Context
}

func GetAccountFromStorage(s *server, name string) (models.Account, error) {
	row := s.db.QueryRowContext(s.ctx, "SELECT name, amount FROM accounts WHERE name=$1", name)

	account := models.Account{}
	err := row.Scan(&account.Name, &account.Amount)

	switch {
	case err == sql.ErrNoRows:
		return models.Account{}, fmt.Errorf("account not found")
	case err != nil:

		return models.Account{}, fmt.Errorf("failed to get account: %w", err)
	default:
		return account, nil
	}
}

func (s *server) Get(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountReply, error) {
	account, err := GetAccountFromStorage(s, req.GetName())
	if err != nil {
		return nil, err
	}
	return &proto.GetAccountReply{Name: account.Name, Amount: int32(account.Amount)}, nil
}

func (s *server) Create(ctx context.Context, req *proto.CreateAccountRequest) (*proto.Empty, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}

	_, err := GetAccountFromStorage(s, req.GetName())
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "account already exists")
	}
	_, err = s.db.ExecContext(ctx, "INSERT INTO accounts(name, amount) VALUES($1, $2)", req.GetName(), req.GetAmount())
	if err != nil {
		return nil, fmt.Errorf("failed to insert account: %w", err)
	}
	return &proto.Empty{}, nil
}

func (s *server) ChangeAmount(ctx context.Context, req *proto.PatchAccountRequest) (*proto.Empty, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	_, err := GetAccountFromStorage(s, req.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to found account: %w", err)
	}
	_, err = s.db.ExecContext(ctx, "UPDATE accounts SET amount = $1 WHERE name = $2", req.GetAmount(), req.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to change amount: %w", err)
	}
	return &proto.Empty{}, nil
}

func (s *server) ChangeName(ctx context.Context, req *proto.ChangeAccountRequest) (*proto.Empty, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	if len(req.GetNewName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty new name")
	}

	_, err := GetAccountFromStorage(s, req.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to found account: %w", err)
	}
	_, err = GetAccountFromStorage(s, req.GetNewName())
	if err == nil {
		return nil, fmt.Errorf("account with new name already exists")
	}
	_, err = s.db.ExecContext(ctx, "UPDATE accounts SET name = $1 WHERE name = $2", req.GetNewName(), req.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to change name: %w", err)
	}
	return &proto.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.Empty, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	_, err := GetAccountFromStorage(s, req.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to found account: %w", err)
	}
	_, err = s.db.ExecContext(ctx, "DELETE FROM accounts WHERE name=$1", req.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to change name: %w", err)
	}
	return &proto.Empty{}, nil
}

func main() {
	connectionString := "host=0.0.0.0 port=5432 dbname=postgres user=postgres password=mysecretpassword"
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	defer db.Close()

	ctx := context.Background()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4567))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterAccountServer(s, New(db, ctx))
	if err := s.Serve(lis); err != nil {
		panic(err)
	}

}
