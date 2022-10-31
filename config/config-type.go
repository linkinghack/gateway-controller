package config

type GlobalConfig struct {
	ServerConfig    ControlPlaneAPIServerConfig `json:"serverConfig" yaml:"serverConfig"`
	LogConfig       LoggerConfig                `json:"logConfig" yaml:"logConfig"`
	DBConfig        DBConfig                    `json:"dbConfig" yaml:"dbConfig"`
	EtcdConfig      EtcdConfig                  `json:"etcdConfig" yaml:"etcdConfig"`
	XdsServerConfig XdsServerConfig             `json:"xdsServerConfig" yaml:"xdsServerConfig"`
}

type DBConfig struct {
	Host                 string            `json:"host" yaml:"host"`
	Port                 int               `json:"port" yaml:"port"`
	Database             string            `json:"database" yaml:"database"`
	ConnectionParameters map[string]string `json:"connectionParameters" yaml:"connectionParameters"`
	User                 string            `json:"user" yaml:"user"`
	Password             string            `json:"password" yaml:"password"`
	EngineType           string            `json:"engineType" yaml:"engineType"`         // mysql, postgres
	TargetDatabase       string            `json:"targetDatabase" yaml:"targetDatabase"` // dbname
	TablesPrefix         string            `json:"tablesPrefix" yaml:"tablesPrefix"`
	MaxPoolSize          int               `json:"maxPoolSize" yaml:"maxPoolSize"`
	MaxIdleSize          int               `json:"maxIdleSize" yaml:"maxIdleSize"`
}

type EtcdConfig struct {
	Endpoints    []string `json:"endpoints" yaml:"endpoints"`
	AuthType     string   `json:"authType" yaml:"authType"`
	Username     string   `json:"username" yaml:"username"`
	Password     string   `json:"password" yaml:"password"`
	KeyNamespace string   `json:"keyNamespace" yaml:"keyNamespace"` // Auto prepend key prefix, like: gw-control/xds (without trailing slash, it will be added automatically)
}

const (
	EtcdAuthTypeNone     = "None"
	EtcdAuthTypePassowrd = "Password"
)

type LoggerConfig struct {
	Level     string `json:"level" yaml:"level"`
	LogFormat string `json:"logFormat" yaml:"logFormat"` // text, json
	Color     bool   `json:"color" yaml:"color"`         // console log color
	FileLog   struct {
		Enable     bool   `json:"enable" yaml:"enable"`
		LogFileDir string `json:"logFileDir" yaml:"logFileDir"`
	} `json:"fileLog" yaml:"fileLog"`
	//GrayLog struct {
	//	Enable               bool
	//	GrayLogServerAddress string
	//}
}

const (
	LogLevelWarn  = "warn"
	LogLevelInfo  = "info"
	LogLevelDebug = "debug"
	LogLevelTrace = "trace"

	LogFormatText = "text"
	LogFormatJson = "json"
)

type ControlPlaneAPIServerConfig struct {
	ListenAddr  string   `json:"listenAddr" yaml:"listenAddr"`
	GinMode     string   `json:"ginMode" yaml:"ginMode"` // debug, release, test
	CorsOrigins []string `json:"corsOrigins" yaml:"corsOrigins"`
}

type XdsServerConfig struct {
	ListenAddr                  string `json:"listenAddr" yaml:"listenAddr"`
	GrpcKeepaliveSeconds        uint64 `json:"grpcKeepaliceSeconds" yaml:"grpcKepaliveSeconds"`
	GrpcKeepaliveTimeoutSeconds uint64 `json:"grpcKeepaliveTimeoutSeconds" yaml:"grpcKeepaliveTimeoutSeconds"`
	GrpcKeepaliveMinTimeSeconds uint64 `json:"grpcKeepaliveMinTimeSeoncds" yaml:"grpcKeepaliveMinTimeSeoncds"`
	GrpcMaxConcurrentStreams    uint64 `json:"grpcMaxConcurrentStreams" yaml:"grpcMaxConcurrentStreams"`
}
