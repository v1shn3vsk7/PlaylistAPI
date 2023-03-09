package server

import (
	"context"
	"fmt"
	"github.com/v1shn3vsk7/PlaylistAPI/internal/database/postgres"
	pb "github.com/v1shn3vsk7/PlaylistAPI/internal/server/grpc/proto"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/playlist"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/song"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	playlist *playlist.Playlist
	db       *postgres.Postgres
	server   *grpc.Server
	pb.UnimplementedPlayerServer
}

func NewServer(playlist *playlist.Playlist, server *grpc.Server, db *postgres.Postgres) *Server {
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

	postgres, err := postgres.NewDb(os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	defer postgres.Close()

	p := playlist.New()

	s := NewServer(p, grpcServer, postgres)

	songs, err := s.db.GetSongs()
	if err != nil {
		return err
	}
	for _, s := range songs {
		log.Printf("Loaded '%s' by '%s' from db\n", s.Name, s.Artist)
		p.AddSong(s)
	}
	if err := s.serve(p, &listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) serve(playlist *playlist.Playlist, lis *net.Listener) error {
	pb.RegisterPlayerServer(s.server, &Server{playlist: playlist, db: s.db})
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
	newSong := &song.Song{
		Name:     in.Name,
		Artist:   in.Artist,
		Duration: time.Duration(in.Duration) * time.Second,
	}

	if err := s.db.AddSong(newSong); err != nil {
		err = status.Error(codes.Internal, err.Error())
		return &pb.Response{}, err
	}
	s.playlist.AddSong(newSong)

	return &pb.Response{Result: fmt.Sprintf("Song '%s' added to the playlist", newSong.Name)}, nil
}

func (s *Server) Edit(ctx context.Context, in *pb.EditRequest) (*pb.Response, error) {
	prevSong := &song.Song{
		Name:   in.PrevName,
		Artist: in.PrevArtist,
	}

	if s.playlist.IsPlaying && s.playlist.IsCurrentSong(prevSong) {
		err := status.Error(codes.InvalidArgument, "Can not delete playing song")
		return &pb.Response{}, err
	}

	newSong := &song.Song{
		Name:     in.NewName,
		Artist:   in.NewArtist,
		Duration: time.Duration(in.NewDuration) * time.Second,
	}

	id, err := s.db.FindSong(prevSong)
	if err != nil {
		err = status.Error(codes.NotFound, err.Error())
		return &pb.Response{}, err
	}

	s.db.EditSong(newSong, id)
	s.playlist.Edit(prevSong, newSong)

	return &pb.Response{Result: "Successfully edited song"}, nil
}

func (s *Server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Response, error) {
	song := &song.Song{
		Name:   in.Name,
		Artist: in.Artist,
	}

	if s.playlist.IsPlaying && s.playlist.IsCurrentSong(song) {
		err := status.Error(codes.InvalidArgument, "Can not delete playing song")
		return &pb.Response{}, err
	}

	if err := s.db.DeleteSong(song); err != nil {
		err = status.Error(codes.NotFound, err.Error())
		return &pb.Response{}, err
	}

	s.playlist.Delete(song)

	return &pb.Response{Result: fmt.Sprintf("Song '%s' by '%s' was deleted",
		song.Name, song.Artist)}, nil
}

func (s *Server) Status(ctx context.Context, in *emptypb.Empty) (*pb.Response, error) {
	return &pb.Response{Result: s.playlist.Status()}, nil
}
