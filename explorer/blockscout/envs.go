package blockscout

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func patchDotEnv(path string, envs map[string]string) error {
	err := restoreFromBackup(path)
	if err != nil {
		return err
	}

	dotEnv, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dotEnv.Close()

	env, err := godotenv.Parse(dotEnv)
	if err != nil {
		return err
	}

	env = mergeMaps(env, envs)
	return godotenv.Write(env, path)
}

func restoreFromBackup(path string) error {
	backupPath := path + ".bak"
	_, err := os.Stat(backupPath)
	if os.IsNotExist(err) {
		return copyFile(path, backupPath)
	} else if err != nil {
		return err
	} else {
		return copyFile(backupPath, path)
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func parseEnvs(envs []string) (map[string]string, error) {
	envMap := make(map[string]string)
	for _, env := range envs {
		kv := strings.Split(env, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid env: %s", env)
		}
		envMap[kv[0]] = kv[1]
	}
	return envMap, nil
}

func mergeMaps(maps ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}
