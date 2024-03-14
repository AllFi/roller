package blockscout

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dymensionxyz/roller/config"
	"github.com/go-git/go-git/v5"
	"github.com/joho/godotenv"
)

const (
	Repository    = "https://github.com/blockscout/blockscout.git"
	DockerCompose = "docker-compose"
)

type Blockscout struct {
	home string
}

func New(home string) *Blockscout {
	return &Blockscout{
		home: fmt.Sprintf("%s/blockscout", home),
	}
}

func (b *Blockscout) Start() error {
	err := b.cloneRepo()
	if err != nil {
		return err
	}

	err = b.configure()
	if err != nil {
		return err
	}

	return b.start()
}

func (b *Blockscout) Stop() error {
	return b.stop()
}

func (b *Blockscout) Clear() error {
	err := os.RemoveAll(filepath.Join(b.home))
	if err != nil {
		fmt.Printf("Error clearing Blockscout data: %s\n", err.Error())
		return err
	}
	return nil
}

func (b *Blockscout) cloneRepo() error {
	fmt.Printf("Cloning Blockscout repository to %s\n", b.home)
	_, err := git.PlainClone(b.home, false, &git.CloneOptions{
		URL:          Repository,
		Progress:     os.Stdout,
		SingleBranch: true,
		Depth:        1,
	})
	if err == git.ErrRepositoryAlreadyExists {
		fmt.Printf("Blockscout repository already exists at %s. Skipping...\n", b.home)
		return nil
	}
	return err
}

func (b *Blockscout) configure() error {
	rollappConfig, err := config.LoadConfigFromTOML(filepath.Dir(b.home))
	if err != nil {
		return err
	}

	backendPath, err := b.backendDotEnvPath()
	if err != nil {
		return err
	}

	env := make(map[string]string)
	env["NETWORK"] = rollappConfig.RollappID
	env["SUBNETWORK"] = rollappConfig.RollappID
	env["COIN"] = rollappConfig.Denom

	err = b.patchDotEnv(backendPath, env)
	if err != nil {
		return err
	}

	frontendPath, err := b.frontendDotEnvPath()
	if err != nil {
		return err
	}

	env = make(map[string]string)
	env["NEXT_PUBLIC_NETWORK_NAME"] = rollappConfig.RollappID
	env["NEXT_PUBLIC_NETWORK_SHORT_NAME"] = rollappConfig.RollappID
	env["NEXT_PUBLIC_NETWORK_CURRENCY_NAME"] = rollappConfig.Denom
	env["NEXT_PUBLIC_NETWORK_CURRENCY_SYMBOL"] = rollappConfig.Denom

	return b.patchDotEnv(frontendPath, env)
}

func (b *Blockscout) start() error {
	fmt.Printf("Starting Blockscout...\n")
	return b.dockerComposeExecute("up", "-d")
}

func (b *Blockscout) stop() error {
	fmt.Printf("Stopping Blockscout...\n")
	return b.dockerComposeExecute("down", "-v")
}

func (b *Blockscout) dockerComposeExecute(commandArgs ...string) error {
	yamlPath, err := b.dockerComposeYamlPath()
	if err != nil {
		return err
	}
	args := []string{
		"--file", yamlPath,
	}
	args = append(args, commandArgs...)
	cmd := exec.Command(
		DockerCompose, args...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (b *Blockscout) patchDotEnv(path string, envs map[string]string) error {
	dotEnv, err := os.Open(path)
	if err != nil {
		return err
	}

	env, err := godotenv.Parse(dotEnv)
	if err != nil {
		return err
	}

	for k, v := range envs {
		env[k] = v
	}

	return godotenv.Write(env, path)
}

func (b *Blockscout) dockerComposeYamlPath() (string, error) {
	return filepath.Abs(fmt.Sprintf("%s/docker-compose/docker-compose.yml", b.home))
}

func (b *Blockscout) backendDotEnvPath() (string, error) {
	return filepath.Abs(fmt.Sprintf("%s/docker-compose/envs/common-blockscout.env", b.home))
}

func (b *Blockscout) frontendDotEnvPath() (string, error) {
	return filepath.Abs(fmt.Sprintf("%s/docker-compose/envs/common-frontend.env", b.home))
}
