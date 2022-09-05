package conf

import (
	"os"
	"path/filepath"
	"server/internal/api/metrics"
	"server/internal/utils"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jkstack/jkframe/conf/kvconf"
	"github.com/jkstack/jkframe/logging"
	runtime "github.com/jkstack/jkframe/utils"
)

const (
	defaultCacheDir       = "/opt/smartagent-server/cache"
	defaultCacheThreshold = 80
	defaultDataDir        = "/opt/smartagent-server/data"
	defaultPluginDir      = "/opt/smartagent-server/plugins"
	defaultLogDir         = "/opt/smartagent-server/logs"
	defaultLogSize        = utils.Bytes(50 * 1024 * 1024)
	defaultLogRotate      = 7
	defaultConnectLimit   = 100
)

type Configure struct {
	Listen         uint16      `kv:"listen"`
	CacheDir       string      `kv:"cache_dir"`
	CacheThreshold uint        `kv:"cache_threshold"`
	LogDir         string      `kv:"log_dir"`
	LogSize        utils.Bytes `kv:"log_size"`
	LogRotate      int         `kv:"log_rotate"`
	ConnectLimit   int         `kv:"connect_limit"`
	Metrics        struct {
		Kafka string `kv:"kafka_addr"`
		Topic string `kv:"kafka_topic"`
	} `kv:"metrics"`
	// runtime
	WorkDir    string
	MetricsCli sarama.AsyncProducer
}

func Load(dir, abs string) *Configure {
	f, err := os.Open(dir)
	runtime.Assert(err)
	defer f.Close()

	var ret Configure
	runtime.Assert(kvconf.NewDecoder(f).Decode(&ret))
	ret.check(abs)

	ret.WorkDir, _ = os.Getwd()

	if len(ret.Metrics.Kafka) > 0 {
		cfg := sarama.NewConfig()
		cfg.Producer.Flush.Bytes = 1024 * 1024 // 1MB
		cfg.Producer.Flush.Messages = 100
		cfg.Producer.Flush.Frequency = time.Second
		ret.MetricsCli, err = sarama.NewAsyncProducer([]string{ret.Metrics.Kafka}, cfg)
		runtime.Assert(err)
		go metrics.HandleReportError(ret.MetricsCli)
	}

	return &ret
}

func (cfg *Configure) check(abs string) {
	if cfg.Listen == 0 {
		panic("invalid listen config")
	}
	if len(cfg.CacheDir) == 0 {
		logging.Info("reset conf.cache_dir to default path: %s", defaultCacheDir)
		cfg.CacheDir = defaultCacheDir
	} else if !filepath.IsAbs(cfg.CacheDir) {
		cfg.CacheDir = filepath.Join(abs, cfg.CacheDir)
	}
	if len(cfg.LogDir) == 0 {
		logging.Info("reset conf.log_dir to default path: %s", defaultLogDir)
		cfg.LogDir = defaultLogDir
	} else if !filepath.IsAbs(cfg.LogDir) {
		cfg.LogDir = filepath.Join(abs, cfg.LogDir)
	}
	if cfg.LogSize == 0 {
		logging.Info("reset conf.log_size to default size: %s", defaultLogSize.String())
		cfg.LogSize = defaultLogSize
	}
	if cfg.LogRotate == 0 {
		logging.Info("reset conf.log_roate to default count: %d", defaultLogRotate)
		cfg.LogRotate = defaultLogRotate
	}
	if cfg.CacheThreshold == 0 || cfg.CacheThreshold > 100 {
		logging.Info("reset conf.cache_threshold to default limit: %d", defaultCacheThreshold)
		cfg.CacheThreshold = defaultCacheThreshold
	}
	if cfg.ConnectLimit == 0 {
		logging.Info("reset conf.connect_limit to default limit: %d", defaultConnectLimit)
		cfg.ConnectLimit = defaultConnectLimit
	}
}
