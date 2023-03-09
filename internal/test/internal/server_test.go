package test_internal

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/v1shn3vsk7/PlaylistAPI/internal/server"
	pb "github.com/v1shn3vsk7/PlaylistAPI/internal/server/grpc/proto"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/playlist"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := playlist.New()
	srv := grpc.NewServer()
	s := server.NewServer(p, srv, db)
	pb.RegisterPlayerServer(srv, s)

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when starting a server", err)
	}

	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPlayerClient(conn)

	_, err = client.Play(context.Background(), &emptypb.Empty{})
	if err.Error() != "rpc error: code = FailedPrecondition desc = playlist is empty" {
		t.Fatalf("Play() got unexpected error")
	}

	//addReq := &pb.AddRequest{
	//	Name:     "Song-1",
	//	Artist:   "Artist-1",
	//	Duration: 400,
	//}
	//
	//mock.ExpectBegin()
	//mock.ExpectExec("INSERT INTO songs VALUES").WithArgs(addReq.Name, addReq.Artist, addReq.Duration).WillReturnResult(sqlmock.NewResult(1, 1))
	//mock.ExpectRollback()
	//
	//_, err = client.Add(context.Background(), addReq)
	//if err != nil {
	//	t.Fatalf("Add() got unexpected error: %v", err)
	//}

}
