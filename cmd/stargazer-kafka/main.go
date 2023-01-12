package main

import (
	"context"
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/config"
	"github.com/entiros/stargazer-kafka/internal/log"
	"github.com/entiros/stargazer-kafka/internal/metrics"
	"github.com/entiros/stargazer-kafka/internal/system"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var DefaultHealthPort = 8081
var DefaultMetricsPort = 9090

func makeDir(dir string) error {

	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		err := os.Mkdir(dir, 0777)
		if err != nil {
			return err
		}
		return nil
	}
	return err

}

func main() {

	log.Logger.Debugf("Starting Stargazer")
	defer log.Logger.Debugf("Stargazer closing down.")

	if len(os.Args) < 2 {
		log.Logger.Fatal("Start with configuration file name or name of directory with multiple configuration files")
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	_, err := os.Stat(os.Args[1])
	if err != nil {
		log.Logger.Fatal(err)
	}
	fileName := os.Args[1]

	log.Logger.Debugf("Using config from: %s", fileName)

	srvGroup, srvContext := errgroup.WithContext(ctx)

	srvGroup.Go(func() error {
		return startHealthServer(srvContext, healthPort())
	})

	srvGroup.Go(func() error {
		return runMetricsServer(srvContext, metricsPort())
	})

	srvGroup.Go(func() error {
		return runSync(srvContext, fileName)
	})

	log.Logger.Debugf("Stargazer running.")
	shutdownErr := srvGroup.Wait()

	if shutdownErr != nil && shutdownErr != http.ErrServerClosed {
		log.Logger.Errorf("Shutdown with error: %v", shutdownErr)
	} else {
		log.Logger.Info("Clean shutdown.")
	}

}

func runSync(ctx context.Context, fileName string) error {

loop:
	for {
		next, hasNext := getSystems(fileName, ctx)

		for hasNext() {
			sys, err := next()
			if err != nil {
				metrics.ErrCount.Add(1)
				log.Logger.Errorf("Failed to sync. %v", err)
				continue
			}

			err = sys.PingStarlify(ctx)
			if err != nil {
				metrics.ErrCount.Add(1)
				log.Logger.Errorf("failed to ping starlify. %v", err)
			}
			err = sys.SyncTopics(ctx)
			if err != nil {
				metrics.ErrCount.Add(1)
				log.Logger.Errorf("failed to sync topics for %s, %v ", sys.Name(), err)
			}

			metrics.SyncCount.Add(1)
			time.Sleep(3 * time.Second)

		}

		select {
		case <-ctx.Done():
			break loop
		case <-time.After(10 * time.Second):
		}
		time.Sleep(10 * time.Second)
	}

	return nil
}

func healthPort() int {

	port := os.Getenv("HEALTH_PORT")
	if port != "" {
		port, err := strconv.Atoi(port)
		if err != nil {
			return port
		}
		return DefaultHealthPort
	}
	return DefaultHealthPort
}

func metricsPort() int {

	port := os.Getenv("METRICS_PORT")
	if port != "" {
		port, err := strconv.Atoi(port)
		if err != nil {
			return port
		}
		return DefaultMetricsPort
	}
	return DefaultMetricsPort
}

func errorFoo(msg error) (next func() (*system.System, error), hasNext func() bool) {
	return func() (*system.System, error) {
			return nil, fmt.Errorf("%v", msg)
		}, func() bool {
			return false
		}

}

func getSystems(fileName string, ctx context.Context) (next func() (*system.System, error), hasNext func() bool) {

	info, err := os.Stat(fileName)
	if err != nil {
		return errorFoo(err)
	}
	if info.IsDir() {
		return GetSystems(ctx, fileName)
	} else {
		return GetSystem(ctx, fileName)
	}
}

func GetSystem(ctx context.Context, dir string) (next func() (*system.System, error), hasNext func() bool) {

	var i int
	return func() (*system.System, error) {
			s, err := system.NewSystem(ctx, dir)
			i++
			return s, err
		}, func() bool {
			return i < 1
		}

}

func GetSystems(ctx context.Context, dir string) (next func() (*system.System, error), hasNext func() bool) {

	configFiles, err := config.GetConfigs(dir)
	if err != nil {
		log.Logger.Errorf("Failed to get config files: %v", err)
		return nil, func() bool {
			return false
		}
	}

	var i int

	next = func() (*system.System, error) {
		sys, err := system.NewSystem(ctx, configFiles[i])
		i++
		return sys, err
	}

	hasNext = func() bool {
		return i < len(configFiles)
	}

	return
}

func runMetricsServer(ctx context.Context, metricsPort int) error {

	metricsRouter := gin.New()
	metricsRouter.Use(gin.Recovery())
	metricsRouter.Use(rateLimiter(rate.NewLimiter(5.0, 20)))
	metricsRouter.GET("/metrics", metrics.PrometheusHandler())
	metricsRouter.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/metrics/"))

	metricsSrv := &http.Server{
		Addr:    ":" + strconv.Itoa(metricsPort),
		Handler: metricsRouter,
	}

	return runServer(ctx, metricsSrv)
}

func ready() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func alive() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func startHealthServer(ctx context.Context, healthPort int) error {

	healthRouter := gin.New()
	healthRouter.Use(gin.Recovery())
	healthRouter.Use(rateLimiter(rate.NewLimiter(3.0, 1)))
	healthRouter.GET("/readyz", ready())
	healthRouter.GET("/livez", alive())
	healthRouter.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/readyz", "/livez"))

	healthSrv := &http.Server{
		Addr:    ":" + strconv.Itoa(healthPort),
		Handler: healthRouter,
	}

	return runServer(ctx, healthSrv)

}

func runServer(ctx context.Context, srv *http.Server) error {
	grp, grpCtx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Logger.Errorf("Failed to listen: %v", err)
		}
		return err
	})

	grp.Go(func() error {
		<-grpCtx.Done()
		shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		log.Logger.Debugf("Shutting down server: %s", srv.Addr)
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Logger.Infof("Failed to shutdown gracefully: %v", err)
		}
		log.Logger.Debugf("Shut down gracefully: %s", srv.Addr)
		return err
	})

	return grp.Wait()
}

func rateLimiter(limiter *rate.Limiter) func(c *gin.Context) {

	return func(c *gin.Context) {

		deadline, cancel := context.WithDeadline(c.Request.Context(), time.Now().Add(time.Second*5))
		defer cancel()

		err := limiter.Wait(deadline)
		if err != nil {
			c.AbortWithStatus(http.StatusTooManyRequests)
		}
	}

}
