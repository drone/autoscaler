package engine

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/drone/autoscaler"

	"github.com/rs/zerolog/log"
)

type zombie struct {
	wg sync.WaitGroup

	servers  autoscaler.ServerStore
	provider autoscaler.Provider
	client   clientFunc
	minAge   time.Duration
}

func (z *zombie) DetectZombies(ctx context.Context) error {
	logger := log.Ctx(ctx)
	logger.Debug().Msg("detect zombies called")

	serverList, err := z.servers.List(ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Error occured while fetching server list")
		return nil
	}

	// we filter out all "old" stopped agents. these remain in the database and should not be zombie checked
	var serverListExcludingStopped []*autoscaler.Server
	for _, server := range serverList {
		if server.State != autoscaler.StateStopped {
			serverListExcludingStopped = append(serverListExcludingStopped, server)
		}
	}

	logger.Debug().
		Int("server count", len(serverListExcludingStopped)).
		Msg("detect zombies server count")

	for _, server := range serverListExcludingStopped {
		z.wg.Add(1)
		go func(server *autoscaler.Server) {
			z.detectZombieAndDelete(ctx, server)
			z.wg.Done()
		}(server)
	}
	return nil
}

func (z *zombie) detectZombieAndDelete(ctx context.Context, instance *autoscaler.Server) error {
	logger := log.Ctx(ctx).With().
		Str("ip", instance.Address).
		Str("name", instance.Name).
		Logger()

	logger.Debug().Msg("detect zombie called")

	conn, err := net.DialTimeout("tcp", instance.Address+":2376", time.Duration(5*time.Second))
	if err == nil {
		logger.Debug().Msg("instance was found alive. no zombie")
		conn.Close()
		return nil
	}

	// Ignore Pending. These servers do not have a "Created" field yet so we cannot (yet) assert if they are of age.
	// Ignore Shutdown. If a zombie is detected we set the state to Shutdown. It's already scheduled for deletion
	// Ignore Stopping. A termination operation is in progress for this agent
	if instance.State == autoscaler.StatePending || instance.State == autoscaler.StateShutdown ||
		instance.State == autoscaler.StateStopping {
		logger.Debug().
			Msg("instance in state Pening/Shutdown/Stopping. not viable for for zombie check")
		return nil
	}

	// If the code arrives here we detected a Zombie
	// check if the agent is older than :minAge: (see config. default 5 minutes). if so schedule for removal
	if time.Now().After(time.Unix(instance.Created, 0).Add(z.minAge)) {
		instance.State = autoscaler.StateShutdown
		err := z.servers.Update(ctx, instance)
		logger.Info().
			Msg("zombie detected. Setting to state shutdown so it will get deleted")
		if err != nil {
			logger.Error().
				Err(err).
				Str("state", "shutdown").
				Msg("cannot update server state")
			return nil
		}
	} else {
		logger.Debug().
			Msg("server is too young to be declared a zombie. doing nothing")
	}
	return nil
}
