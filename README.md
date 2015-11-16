# Usage

## Config

database.staging.yml 
```
port: 8000
log:
  filename: charity
  maxdays: 30
```

config.go
```
var Config = struct {
    Env  string `env:"ENV" default:"local"`
    Port uint   `env:"PORT" default:"7000"`
    Log struct {
        FileName string
        Maxdays  int `default:"30"`
    }
}{}

func init() {
    if err := configor.Load(&Config, "config/database.yml"); err != nil {
        panic(err)
    }
}
```

```
import "github.com/qor/log"

router := gin.New()
router.Use(log.Logger(config.Config.Log.FileName,config.Config.Log.Maxdays))
```