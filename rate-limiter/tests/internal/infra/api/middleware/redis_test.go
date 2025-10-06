package middleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	redisInternal "github.com/deal-machine/go-expert/rate-limiter/internal/infra/persistence/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	testcontainers.Container
	URI string
}

var redisContainer *RedisContainer
var redisClient *redis.Client

func TestMain(m *testing.M) {
	ctx := context.Background()
	t := &testing.T{}

	redisContainer = SetupRedisContainer(ctx, t)
	client, err := redisInternal.NewRedisClient(ctx, redisContainer.URI)
	if err != nil {
		log.Fatal(err)
	}
	redisClient = client

	code := m.Run()

	client.Close()
	redisContainer.Terminate(ctx)
	os.Exit(code)
}

func SetupRedisContainer(ctx context.Context, t *testing.T) *RedisContainer {
	t.Helper()
	req := testcontainers.ContainerRequest{
		Image:        "redis",
		ExposedPorts: []string{"6379"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}
	genericContainer := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}
	container, err := testcontainers.GenericContainer(ctx, genericContainer)
	require.NoError(t, err)
	mappedPort, err := container.MappedPort(ctx, "6379")
	require.NoError(t, err)
	hostIP, err := container.Host(ctx)
	require.NoError(t, err)
	uri := fmt.Sprintf("%s:%s", hostIP, mappedPort.Port())
	redisContainer := &RedisContainer{
		URI:       uri,
		Container: container,
	}
	require.NotNil(t, redisContainer.URI)
	require.NotNil(t, redisContainer.Container)
	return redisContainer
}
