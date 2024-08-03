package adapters

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type (
	Config map[string]string
)

func (c *Config) Load(filePath string) {
	mp, err := godotenv.Read(filePath)
	if err != nil {
		logger.Println("config: load from host env")
		mp = make(map[string]string)
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			mp[pair[0]] = pair[1]
		}
	}
	*c = mp
}
