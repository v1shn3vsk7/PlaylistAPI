package test_internal

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	mockDb "github.com/v1shn3vsk7/PlaylistAPI/internal/database/mock"
	"github.com/v1shn3vsk7/PlaylistAPI/internal/server"
	pb "github.com/v1shn3vsk7/PlaylistAPI/internal/server/grpc/proto"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/playlist"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	DB := mockDb.MockDb{db}
	p := playlist.New()
	srv := grpc.NewServer()
	s := server.NewServer(p, srv, &DB)
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

	addReq := &pb.AddRequest{
		Name:     "Song-1",
		Artist:   "Artist-1",
		Duration: 400,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO songs").WithArgs(addReq.Name, addReq.Artist, addReq.Duration).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err = client.Add(context.Background(), addReq)
	if err != nil {
		t.Fatalf("Add() got unexpected error: %v", err)
	}

	_, err = client.Play(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("Add() got unexpected error: %v", err)
	}

	_, err = client.Pause(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("Pause() got unexpected error: %v", err)
	}

	editReq := &pb.EditRequest{
		PrevName:    "Song-1",
		PrevArtist:  "Artist-1",
		NewName:     "Song-2",
		NewArtist:   "Artist-2",
		NewDuration: 500,
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT id FROM songs").WithArgs(editReq.PrevName, editReq.PrevArtist)
	mock.ExpectExec("UPDATE songs").WithArgs(editReq.NewName, editReq.NewArtist, editReq.NewDuration)
	mock.ExpectCommit()

	_, err = client.Edit(context.Background(), editReq)
	if err != nil {
		t.Fatalf("Edit() got unexpected error: %v", err)
	}

	_, err = client.Play(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("Add() got unexpected error: %v", err)
	}

}
