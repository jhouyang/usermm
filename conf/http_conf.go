package conf

// HTTPConf conf object for httpserver
type HTTPConf struct {
    Server struct {
        Port int    `yaml:"port"`
        IP   string `yaml:"ip"`
    }
    Image struct {
        Prefixurl string `yaml:"prefixurl"`
        Savepath  string `yaml:"savepath"`
        Maxsize   int    `yaml:"maxsize"`
    }
    Logic struct {
        Tokenexpire int `yaml:"tokenexpire"`
    }
    Log struct {
        Logfile  string `yaml:"logfile"`
        Loglevel string `yaml:"loglevel"`
        Maxdays  string `yaml:"maxdays"`
    }
    Rpcserver struct {
        Addr string `yaml:"addr"`
    }
    Pool struct {
        Initsize   uint32 `yaml:"initsize"`
        Capacity   uint32 `yaml:"capacity"`
        Maxidle    uint8  `yaml:"maxidle"`
        Gettimeout uint8  `yaml:"gettimeout"`
    }
}
