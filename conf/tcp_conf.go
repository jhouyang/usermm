package conf

// TCPConf go object to tcpserver.yaml
type TCPConf struct {
    Server struct {
        Port int `yaml:"port"`
    }
    Log struct {
        Logfile  string `yaml:"logfile"`
        Loglevel string `yaml:"loglevel"`
        Maxdays  string `yaml:"maxdays"`
    }
    Db struct {
        Host   string `yaml:"host"`
        User   string `yaml:"user"`
        Passwd string `yaml:"passwd"`
        Db     string `yaml:"db"`
        Conn struct {
            Maxidle  int `yaml:"maxidle"`
            Maxopen int `yaml:"maxopen"`
        }
    }
    Redis struct {
        Addr     string `yaml:"addr"`
        Db       int    `yaml:"db"`
        Passwd   string `yaml:"passwd"`
	Poolsize int    `yaml:"poolsize"`
        Cache struct {
            Tokenexpired int  `yaml:"tokenexpired"`
            Userexpired  int  `yaml:"userexpired"`
        }
    }
}
