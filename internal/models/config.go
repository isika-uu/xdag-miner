package models

// XMRigConfig XMRig配置结构
type XMRigConfig struct {
	API      APIConfig     `json:"api"`
	HTTP     HTTPConfig    `json:"http"`
	Autosave bool          `json:"autosave"`
	CPU      CPUConfig     `json:"cpu"`
	Pools    []PoolConfig  `json:"pools"`
	RandomX  RandomXConfig `json:"randomx"`
	LogFile  *string       `json:"log-file"`
}

// APIConfig API配置
type APIConfig struct {
	ID       *string `json:"id"`
	WorkerID *string `json:"worker-id"`
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Enabled     bool    `json:"enabled"`
	Host        string  `json:"host"`
	Port        int     `json:"port"`
	AccessToken *string `json:"access-token"`
	Restricted  bool    `json:"restricted"`
}

// CPUConfig CPU配置
type CPUConfig struct {
	Enabled        bool `json:"enabled"`
	HugePages      bool `json:"huge-pages"`
	MaxThreadsHint int  `json:"max-threads-hint"`
	Priority       *int `json:"priority"`
	ASM            bool `json:"asm"`
}

// PoolConfig 矿池配置
type PoolConfig struct {
	Algo     *string `json:"algo"`
	Coin     *string `json:"coin"`
	URL      string  `json:"url"`
	User     string  `json:"user"`
	Pass     string  `json:"pass"`
	RigID    *string `json:"rig-id"`
	Nicehash bool    `json:"nicehash"`
	Enabled  bool    `json:"enabled"`
	TLS      bool    `json:"tls"`
}

// RandomXConfig RandomX算法配置
type RandomXConfig struct {
	Init       int    `json:"init"`
	InitAVX2   int    `json:"init-avx2"`
	Mode       string `json:"mode"`
	OneGBPages bool   `json:"1gb-pages"`
	NUMA       bool   `json:"numa"`
}

// MinerStatus 挖矿状态
type MinerStatus struct {
	Running   bool    `json:"running"`
	Hashrate  float64 `json:"hashrate"`
	Threads   int     `json:"threads"`
	Uptime    int64   `json:"uptime"`
	Pool      string  `json:"pool"`
	Algorithm string  `json:"algorithm"`
	Connected bool    `json:"connected"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	CPUModel     string `json:"cpuModel"`
	CPUCores     int    `json:"cpuCores"`
	TotalMemory  uint64 `json:"totalMemory"`
	XMRigVersion string `json:"xmrigVersion"`
}
