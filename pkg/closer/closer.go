package closer

import (
	"context"
	"fmt"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type closer struct {
	mu              sync.Mutex
	funcs           []Func
	shutdownTimeout time.Duration
}

type Closer interface {
	Add(f Func)
	Close(ctx context.Context) error
}

func New(shutdownTimeout time.Duration) *closer {
	return &closer{
		shutdownTimeout: shutdownTimeout,
		mu:              sync.Mutex{},
		funcs:           []Func{},
	}
}

type Func func() error

func (c *closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		msgs     = make([]string, 0, len(c.funcs))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, f := range c.funcs {
			if err := f(); err != nil {
				msgs = append(msgs, err.Error())
			}
		}
		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("closer cancelled: %v", ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf("closer finished with error(s): %s", strings.Join(msgs, "; "))
	}

	return nil
}

func (c *closer) GracefulShutdown() {
	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	<-stopCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), c.shutdownTimeout)
	defer cancel()

	shutdown := make(chan bool)

	go func() {
		err := c.Close(shutdownCtx)
		if err != nil {
			log.Err(err).Msg("closer error")
		}
		shutdown <- true
	}()

	select {
	case <-shutdownCtx.Done():
		log.Err(shutdownCtx.Err()).Msg("server shutdown")
	case <-shutdown:
		log.Info().Msg("shutting down gracefully")
	}
}
