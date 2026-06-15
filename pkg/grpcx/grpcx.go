// Package grpcx — спільні хелпери для підняття gRPC-серверів і клієнтів усіх
// сервісів: реєстрація health-check, graceful shutdown, dial з очікуванням.
package grpcx

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Serve піднімає gRPC-сервер на addr, дає register зареєструвати сервіси,
// вмикає health + reflection і блокується до SIGINT/SIGTERM, після чого робить
// graceful stop. Зручно викликати прямо з main кожного сервісу.
func Serve(addr string, log *slog.Logger, register func(*grpc.Server)) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	register(srv)

	hs := health.NewServer()
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(srv, hs)
	reflection.Register(srv)

	go func() {
		log.Info("gRPC слухає", "addr", addr)
		if err := srv.Serve(lis); err != nil {
			log.Error("gRPC serve впав", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("зупиняюсь...")

	done := make(chan struct{})
	go func() {
		srv.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		srv.Stop()
	}
	return nil
}

// Dial відкриває (lazy) клієнтське з'єднання без TLS — усе всередині приватної
// мережі docker-compose. З'єднання не блокується: gRPC сам переконнектиться.
func Dial(target string) (*grpc.ClientConn, error) {
	return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// WaitHealthy блокується доки сервіс target не відповість SERVING або не вийде
// дедлайн. Корисно для orchestrator на старті, поки провайдери ще піднімаються.
func WaitHealthy(ctx context.Context, target string, log *slog.Logger) {
	conn, err := Dial(target)
	if err != nil {
		log.Warn("не вдалось створити клієнт", "target", target, "err", err)
		return
	}
	defer conn.Close()
	hc := healthpb.NewHealthClient(conn)
	for {
		if ctx.Err() != nil {
			return
		}
		cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		resp, err := hc.Check(cctx, &healthpb.HealthCheckRequest{})
		cancel()
		if err == nil && resp.GetStatus() == healthpb.HealthCheckResponse_SERVING {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}
