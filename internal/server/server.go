package server

import (
	"database/sql"
	"github.com/v1shn3vsk7/PlaylistAPI/internal/database/postgres"
	pb "github.com/v1shn3vsk7/PlaylistAPI/internal/server/grpc/proto"
	"github.com/v1shn3vsk7/PlaylistAPI/pkg/playlist"
	"google.golang.org/grpc"
	"net"
	"os"
)

type Server struct {
	playlist *playlist.Playlist
	db       *sql.DB
	server   *grpc.Server
	pb.UnimplementedPlayerServer
}

func NewServer(playlist *playlist.Playlist, server *grpc.Server) *Server {
	return &Server{
		playlist: playlist,
		server:   server,
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

	songs := postgres.GetSongs(db)

	p := playlist.New()
	for _, s := range songs {
		p.AddSong(&s)
	}

	s := NewServer(p, grpcServer)
	if err := s.Serve(p, &listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) Serve(playlist *playlist.Playlist, lis *net.Listener) error {
	pb.RegisterPlayerServer(s.server, &Server{ playlist: playlist})

	if err := s.server.Serve(*lis); err != nil {
		return err
	}

	return nil
}