package main

import (
	"awesomeProject/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type Command struct {
	Port    int
	Host    string
	Cmd     string
	Name    string
	Amount  int
	NewName string
}

func main() {
	portVal := flag.Int("port", 4567, "server port")
	hostVal := flag.String("host", "0.0.0.0", "server host")
	cmdVal := flag.String("cmd", "", "command to execute")
	nameVal := flag.String("name", "", "name of account")
	amountVal := flag.Int("amount", 0, "amount of account")
	newNameVal := flag.String("new_name", "", "new name of account")
	flag.Parse()

	cmd := Command{
		Port:    *portVal,
		Host:    *hostVal,
		Cmd:     *cmdVal,
		Name:    *nameVal,
		Amount:  *amountVal,
		NewName: *newNameVal,
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cmd.Host, cmd.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = conn.Close()
	}()

	c := proto.NewAccountClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	if err := do(cmd, c, ctx); err != nil {
		panic(err)
	}
	defer cancel()
}

func do(cmd Command, c proto.AccountClient, ctx context.Context) error {
	switch cmd.Cmd {
	case "create":
		if err := create(cmd, c, ctx); err != nil {
			return fmt.Errorf("create account failed: %w", err)
		}
		return nil
	case "get":
		if err := get(cmd, c, ctx); err != nil {
			return fmt.Errorf("get account failed: %w", err)
		}

		return nil
	case "delete":
		if err := delete(cmd, c, ctx); err != nil {
			return fmt.Errorf("delete account failed: %w", err)
		}

		return nil
	case "change_amount":
		if err := change_amount(cmd, c, ctx); err != nil {
			return fmt.Errorf("change amount failed: %w", err)
		}

		return nil
	case "change_name":
		if err := change_name(cmd, c, ctx); err != nil {
			return fmt.Errorf("change name failed: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("unknown command %s", cmd.Cmd)
	}
}

func create(cmd Command, c proto.AccountClient, ctx context.Context) error {
	_, err := c.Create(ctx, &proto.CreateAccountRequest{Name: cmd.Name, Amount: int32(cmd.Amount)})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("account created")
	return nil
}

func get(cmd Command, c proto.AccountClient, ctx context.Context) error {
	r, err := c.Get(ctx, &proto.GetAccountRequest{Name: cmd.Name})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("account found: name: %s, amount: %d", r.GetName(), r.GetAmount())
	return nil
}

func delete(cmd Command, c proto.AccountClient, ctx context.Context) error {
	_, err := c.Delete(ctx, &proto.DeleteAccountRequest{Name: cmd.Name})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("account deleted")
	return nil
}

func change_name(cmd Command, c proto.AccountClient, ctx context.Context) error {
	_, err := c.ChangeName(ctx, &proto.ChangeAccountRequest{Name: cmd.Name, NewName: cmd.NewName})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("account name changed")
	return nil
}

func change_amount(cmd Command, c proto.AccountClient, ctx context.Context) error {
	_, err := c.ChangeAmount(ctx, &proto.PatchAccountRequest{Name: cmd.Name, Amount: int32(cmd.Amount)})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("account amount changed")
	return nil
}
