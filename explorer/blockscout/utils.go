package blockscout

import (
	"os"

	"github.com/joho/godotenv"
)

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
