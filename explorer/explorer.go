package explorer

import (
	"github.com/dymensionxyz/roller/config"
	"github.com/dymensionxyz/roller/explorer/blockscout"
)

type BlockExplorer interface {
	Start(backendEnvs []string, frontendEnvs []string) error
	Stop() error
	Clear() error
}

type Explorer struct {
	explorerType config.ExplorerType
	BlockExplorer
}

func NewExplorer(explorerType config.ExplorerType, home string) *Explorer {
	var explorer BlockExplorer

	switch explorerType {
	case config.Blockscout:
		explorer = blockscout.New(home)
	default:
		panic("Unknown explorer type")
	}

	return &Explorer{
		explorerType:  explorerType,
		BlockExplorer: explorer,
	}
}
