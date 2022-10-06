package types

type Listener struct {
	Name            string   `json:"name" yaml:"name"`
	ListenAddress   string   `json:"listenAddress" yaml:"listenAddress"`
	ListenPort      uint64   `json:"listenPort" yaml:"listenPort"`
	ListenerPlugins []string `json:"listenerPlugins" yaml:"listenerPlugins"`
}

type TlsSniProxy struct {
	PreferredListener string `json:"preferredListener" yaml:"preferredListener"`
	SniHostPattern    string `json:"sniHostPattern" yaml:"sniHostPattern"`
	TargetHostOrIp    string `json:"targetHostOrIp" yaml:"targetHostOrIp"`
	TargetPort        uint64 `json:"targetPort" yaml:"targetPort"`
}

type TcpProxy struct {
	Name             string `json:"name" yaml:"name"`
	SourceListenIp   string `json:"sourceListenIp" yaml:"sourceListenIp"`
	SourceListenPort uint64 `json:"sourceListenPort" yaml:"sourceListenPort"`
	TargetHostOrIp   string `json:"targetHostOrIp" yaml:"targetHostOrPort"`
	TargetPort       uint64 `json:"targetPort" yaml:"targetPort"`
}

type HttpProxy struct {
	Name string `json:"name" yaml:"name"`
}
