package shutdown

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func StartWithGracefulShutdown(ctx context.Context, srv *http.Server, serverPort string) error {
	serverError := make(chan error, 1)

	go func() error {
		fmt.Println("Server running on port ", serverPort)

		err := http.ListenAndServe(":"+serverPort, srv.Handler)
		if err != nil && err != http.ErrServerClosed {
			serverError <- fmt.Errorf("failed to start server %s", err)
		}

		return nil
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		fmt.Println("Shutting down server...")

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		fmt.Println("Server gracefully shutdown")

		return nil
	case err := <-serverError:
		return err
	}
}
