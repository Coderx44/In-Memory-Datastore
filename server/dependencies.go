package server

import (
	"sync"

	"github.com/Coderx44/gg/commands"
)

type dependencies struct {
	dt *commands.Datastore
}

func InitDependencies() (dependencies, error) {
	dataStore := commands.NewDatastore()
	dataStore.Cond = sync.NewCond(dataStore)
	return dependencies{
		dt: dataStore,
	}, nil
}
