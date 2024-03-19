package blockscout

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dymensionxyz/roller/cmd/utils"
	"github.com/dymensionxyz/roller/config"
)

const (
	BlockscoutRepository       = "https://github.com/blockscout/blockscout.git"
	DockerComposeRelativePath  = "docker-compose/docker-compose.yml"
	BackendDotEnvRelativePath  = "docker-compose/envs/common-blockscout.env"
	FrontendDotEnvRelativePath = "docker-compose/envs/common-frontend.env"
	DockerCompose              = "docker-compose"
)

type Blockscout struct {
	home string
}

func New(home string) *Blockscout {
	return &Blockscout{
		home: filepath.Join(home, "blockscout"),
	}
}

func (b *Blockscout) Start(backendParams []string, frontendParams []string) error {
	backendEnvs, err := parseEnvs(backendParams)
	if err != nil {
		return err
	}

	frontendEnvs, err := parseEnvs(frontendParams)
	if err != nil {
		return err
	}

	err = b.cloneRepo()
	if err != nil {
		return err
	}

	err = b.configure(backendEnvs, frontendEnvs)
	if err != nil {
		return err
	}

	fmt.Printf("Starting Blockscout...\n")
	return utils.ExecBashCmd(b.dockerComposeCommand("up", "-d"))
}

func (b *Blockscout) Stop() error {
	fmt.Printf("Stopping Blockscout...\n")
	return utils.ExecBashCmd(b.dockerComposeCommand("down", "-v"))
}

func (b *Blockscout) Clear() error {
	if b.IsRunning() {
		err := b.Stop()
		if err != nil {
			fmt.Printf("Error stopping Blockscout: %s\n", err.Error())
			return err
		}
	}

	err := os.RemoveAll(filepath.Join(b.home))
	if err != nil && !strings.Contains(err.Error(), "No such file or directory") {
		fmt.Printf("Error clearing Blockscout data: %s\n", err.Error())
		return err
	}

	return nil
}

func (b *Blockscout) IsRunning() bool {
	cmd := b.dockerComposeCommand("ps", "-q")
	stdout, err := utils.ExecBashCommandWithStdout(cmd)
	if err != nil {
		return false
	}
	return stdout.String() != ""
}

func (b *Blockscout) cloneRepo() error {
	cmd := exec.Command(
		"git", "clone", BlockscoutRepository, b.home, "--depth", "1",
	)

	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	fmt.Printf("Cloning Blockscout repository...\n")
	err := utils.ExecBashCmd(cmd)
	if err != nil {
		if strings.Contains(stderr.String(), "already exists") {
			fmt.Printf("Blockscout repository already exists at %s.\n", b.home)
			return nil
		}
		return err
	}
	return nil
}

func (b *Blockscout) configure(backendEnvs map[string]string, frontendEnvs map[string]string) error {
	backendDefaultEnvs, frontendDefaultEnvs := b.envsFromRollerConfig()

	backendDotEnvPath := filepath.Join(b.home, BackendDotEnvRelativePath)
	backendEnvs = mergeMaps(backendDefaultEnvs, backendEnvs)
	err := patchDotEnv(backendDotEnvPath, backendEnvs)
	if err != nil {
		return err
	}

	frontendDotEnvPath := filepath.Join(b.home, FrontendDotEnvRelativePath)
	frontendEnvs = mergeMaps(frontendDefaultEnvs, frontendEnvs)
	return patchDotEnv(frontendDotEnvPath, frontendEnvs)
}

func (b *Blockscout) envsFromRollerConfig() (map[string]string, map[string]string) {
	rollappConfig, err := config.LoadConfigFromTOML(filepath.Dir(b.home))
	if err != nil {
		return nil, nil
	}

	backendEnv := map[string]string{
		"NETWORK":    rollappConfig.RollappID,
		"SUBNETWORK": rollappConfig.RollappID,
		"COIN":       rollappConfig.Denom,
	}

	frontendEnv := map[string]string{
		"NEXT_PUBLIC_NETWORK_NAME":            rollappConfig.RollappID,
		"NEXT_PUBLIC_NETWORK_SHORT_NAME":      rollappConfig.RollappID,
		"NEXT_PUBLIC_NETWORK_CURRENCY_NAME":   rollappConfig.Denom,
		"NEXT_PUBLIC_NETWORK_CURRENCY_SYMBOL": rollappConfig.Denom,
	}

	return backendEnv, frontendEnv
}

func (b *Blockscout) dockerComposeCommand(commandArgs ...string) *exec.Cmd {
	dockerComposePath := filepath.Join(b.home, DockerComposeRelativePath)
	args := []string{
		"--file", dockerComposePath,
	}
	args = append(args, commandArgs...)
	cmd := exec.Command(
		DockerCompose, args...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
