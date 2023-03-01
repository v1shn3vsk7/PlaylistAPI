package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/v1shn3vsk7/PlaylistAPI/internal/database/postgres"
	pb "github.com/v1shn3vsk7/PlaylistAPI/internal/server/grpc/proto"
	"github.com/v1shn3vsk7/PlaylistAPI/internal/utils"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/playlist"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/song"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"os"
)

type Server struct {
	playlist *playlist.Playlist
	db       *sql.DB
	server   *grpc.Server
	pb.UnimplementedPlayerServer
}

func newServer(playlist *playlist.Playlist, server *grpc.Server, db *sql.DB) *Server {
	return &Server{
		playlist: playlist,
		server:   server,
		db:       db,
	}
}

func Start() error {
	listener, err := net.Listen("tcp", ":5536")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	db, err := postgres.NewDb(os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	songs, err := postgres.GetSongs(db)
	if err != nil {
		return err
	}

	p := playlist.New()
	for _, s := range songs {
		p.AddSong(&s)
	}

	s := newServer(p, grpcServer, db)
	if err := s.serve(p, &listener); err != nil {
		return err
	}

	return nil
}


func (s *Server) serve(playlist *playlist.Playlist, lis *net.Listener) error {
	pb.RegisterPlayerServer(s.server, &Server{ playlist: playlist, db: s.db})
	if err := s.server.Serve(*lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) Play(ctx context.Context, in *emptypb.Empty) (*pb.Response, error) {
	if err := s.playlist.Play(); err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return &pb.Response{}, err
	}

	return &pb.Response{Result: "Started playing"}, nil
}

func (s *Server) Pause(ctx context.Context, in *emptypb.Empty) (*pb.Response, error) {
	if err := s.playlist.Pause(); err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return &pb.Response{}, err
	}
	return &pb.Response{Result: "Stopped playback"}, nil
}

func (s *Server) Next(ctx context.Context, in *emptypb.Empty) (*pb.Response, error) {
	if err := s.playlist.Next(); err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return &pb.Response{}, err
	}
	return &pb.Response{Result: "Skipped to the next song"}, nil
}

func (s *Server) Prev(ctx context.Context, in *emptypb.Empty) (*pb.Response, error) {
	if err := s.playlist.Prev(); err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return &pb.Response{}, err
	}
	return &pb.Response{Result: "Moved to the previous song"}, nil
}

func (s *Server) Add(ctx context.Context, in *pb.AddRequest) (*pb.Response, error) {
	dur, err := utils.ParseDuration(in.Duration)
	if err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return &pb.Response{}, err
	}

	newSong := &song.Song{
		Name:     in.Name,
		Artist:   in.Artist,
		Duration: dur,
	}

	if err := postgres.AddSong(s.db, newSong); err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return &pb.Response{}, err
	}
	s.playlist.AddSong(newSong)

	return &pb.Response{Result: fmt.Sprintf("Song '%s' added to the playlist", newSong.Name)}, nil
}

func (s *Server) Edit(ctx context.Context, in *pb.EditRequest) (*pb.Response, error) {
	return &pb.Response{Result: ""}, nil
}

func (s *Server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Response, error) {
	return &pb.Response{Result: ""}, nil
}

