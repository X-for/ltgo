import (
	"encoding/json"
	"os"
	"path/filepaht"
)


type ConfigInfo struct {
	cookie string
	language string
	url string
}

type Config struct {
	path string
	cfg ConfigInfo
}


func getConfigPath() (string, error) {
	path = os.UserHomeDir()
	dir, err := os.ReadDir(path + ".ltgo")
	if err != nil {
		log.Fetal(err)
		err := os.Mkdir(".ltgo", 0755)
		if err != nil {
			log.Fetal(err)
			_, err := Create(dir + "config.json")
		}
	}
	file, err := os.ReadFile(dir + "config.json")
	if err != nil {
		log.Fetal(err)
		file, err := Create(file + "config.json")
		if err != nil {
			log.Fetal(err)
		}
	}
	return file, nil
}

func Load(
