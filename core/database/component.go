package database

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"

	"github.com/iotaledger/hive.go/core/app"
	"github.com/iotaledger/hive.go/core/app/pkg/shutdown"
	hivedb "github.com/iotaledger/hive.go/core/database"
	"github.com/iotaledger/hive.go/core/kvstore"
	"github.com/iotaledger/inx-api-core-v1/pkg/daemon"
	"github.com/iotaledger/inx-api-core-v1/pkg/database"
	"github.com/iotaledger/inx-api-core-v1/pkg/database/engine"
)

const (
	DBVersion uint32 = 2
)

func init() {
	CoreComponent = &app.CoreComponent{
		Component: &app.Component{
			Name:     "database",
			DepsFunc: func(cDeps dependencies) { deps = cDeps },
			Params:   params,
			Provide:  provide,
			Run:      run,
		},
	}
}

type dependencies struct {
	dig.In
	Database        *database.Database
	Echo            *echo.Echo
	ShutdownHandler *shutdown.ShutdownHandler
}

var (
	CoreComponent *app.CoreComponent
	deps          dependencies
)

func provide(c *dig.Container) error {

	type storageOut struct {
		dig.Out
		TangleDatabase kvstore.KVStore `name:"tangleDatabase"`
		UTXODatabase   kvstore.KVStore `name:"utxoDatabase"`
	}

	if err := c.Provide(func() (storageOut, error) {
		CoreComponent.LogInfo("Setting up database ...")

		tangleDatabase, err := engine.StoreWithDefaultSettings(ParamsDatabase.Tangle.Path, false, hivedb.EngineAuto, engine.AllowedEnginesStorageAuto...)
		if err != nil {
			return storageOut{}, err
		}

		utxoDatabase, err := engine.StoreWithDefaultSettings(ParamsDatabase.UTXO.Path, false, hivedb.EngineAuto, engine.AllowedEnginesStorageAuto...)
		if err != nil {
			return storageOut{}, err
		}

		return storageOut{
			TangleDatabase: tangleDatabase,
			UTXODatabase:   utxoDatabase,
		}, err
	}); err != nil {
		return err
	}

	type storageDeps struct {
		dig.In
		TangleDatabase kvstore.KVStore `name:"tangleDatabase"`
		UTXODatabase   kvstore.KVStore `name:"utxoDatabase"`
		NetworkID      uint64          `name:"networkId"`
	}

	if err := c.Provide(func(deps storageDeps) (*database.Database, error) {
		store, err := database.New(deps.TangleDatabase, deps.UTXODatabase, deps.NetworkID, ParamsDatabase.Debug)
		if err != nil {
			return nil, err
		}

		store.PrintSnapshotInfo()

		return store, nil
	}); err != nil {
		return err
	}

	return nil
}

func run() error {

	if err := CoreComponent.Daemon().BackgroundWorker("Close database", func(ctx context.Context) {
		<-ctx.Done()

		CoreComponent.LogInfo("Syncing databases to disk ...")
		if err := deps.Database.CloseDatabases(); err != nil {
			CoreComponent.LogPanicf("Syncing databases to disk ... failed: %s", err)
		}
		CoreComponent.LogInfo("Syncing databases to disk ... done")
	}, daemon.PriorityStopDatabase); err != nil {
		CoreComponent.LogPanicf("failed to start worker: %s", err)
	}

	return nil
}
