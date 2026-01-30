package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Service  ServiceConfig  `mapstructure:"service"`
	GRPC     gRPCConfig     `mapstructure:"grpc"`
	REST     RESTConfig     `mapstructure:"rest"`
	MCP      MCPConfig      `mapstructure:"mcp"`
	Database DatabaseConfig `mapstructure:"database"`
	NATS     NATSConfig     `mapstructure:"nats"`
	Voyage   VoyageConfig   `mapstructure:"voyage"`
	Gemini   GeminiConfig   `mapstructure:"gemini"`
	Buffer   BufferConfig   `mapstructure:"buffer"`
	Search   SearchConfig   `mapstructure:"search"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type LoggingConfig struct {
	Level       string `mapstructure:"level"`       // debug, info, warn, error
	Development bool   `mapstructure:"development"` // pretty print for dev
	Encoding    string `mapstructure:"encoding"`    // json or console
}

type ServiceConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Version     string `mapstructure:"version"`
}

type gRPCConfig struct {
	Port              int           `mapstructure:"port"`
	MaxRecvMsgSize    int           `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize    int           `mapstructure:"max_send_msg_size"`
	ConnectionTimeOut time.Duration `mapstructure:"connection_timeout"`
}

type RESTConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	CORSOrigins  []string      `mapstructure:"cors_origins"`
}

type MCPConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Transport string `mapstructure:"transport"` // "stdio" or "sse"
	SSEPort   int    `mapstructure:"sse_port"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type NATSConfig struct {
	URL             string        `mapstructure:"URL"`
	EnableJetStream bool          `mapstructure:"enable_jetstream"`
	MaxReconnect    int           `mapstructure:"max_reconnect"`
	ReconnectWait   time.Duration `mapstructure:"reconnect_wait"`
	StreamPrefix    string        `mapstructure:"stream_prefix"`
	StreamReplicas  int           `mapstructure:"stream_replicas"`
	MaxDeliver      int           `mapstructure:"max_deliver"`
	AckWait         time.Duration `mapstructure:"ack_wait"`
}

type BufferConfig struct {
	TimeThreshold  time.Duration `mapstructure:"time_threshold"`
	TokenThreshold int           `mapstructure:"token_threshold"`
	FlushInterval  time.Duration `mapstructure:"flush_interval"`
}

type GeminiConfig struct {
	APIKey      string        `mapstructure:"api_key"`
	Model       string        `mapstructur:"model"`
	MaxTokens   int           `mapstructure:"max_tokens"`
	Temperature float32       `mapstructure:"temperature"`
	Timeout     time.Duration `mapstructure:"time_duration"`
}

type VoyageConfig struct {
	APIKey     string        `mapstructure:"api_key"`
	Model      string        `mapstructure:"model"`
	Dimensions int           `mapstructure:"dimensions"`
	BatchSize  int           `mapstructure:"batch_size"`
	Timeout    time.Duration `mapstructure:"timeout"`
	BaseURL    string        `mapstructure:"base_url"`
}

type SearchConfig struct {
	DefaultLimit   int     `mapstructure:"default_limit"`
	MaxLimit       int     `mapstructure:"max_limit"`
	SemanticWeight float32 `mapstructure:"semantich_weight"`
	KeywordWeight  float32 `mapstructure:"keyword_weight"`
	MinScore       float32 `mapstructure:"min_score"`
	RRFConstant    int     `mapstructure:"rrf_constant"`
	PPRDamping     float32 `mapstructure:"ppr_damping"`
	PPRMaxIter     int     `mapstructure:"ppr_max_iter"`
}

type SecurityConfig struct {
	JWTSecret      string        `mapstructure:"jwt_secret"`
	JWTExpiry      time.Duration `mapstructure:"jwt_expiry"`
	EncryptionKey  string        `mapstructure:"encryption_key"`
	RateLimitRPS   int           `mapstructure:"rate_limit_rps"`
	RateLimitBurst int           `mapstructure:"rate_limit_burst"`
}

// Load reads config from env vars (SYNCOGNIZE_ prefix) and optional config file
func Load() (*Config, error) {
	// Debug: show current working directory
	wd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", wd)

	// Load .env file - try multiple paths
	err := godotenv.Load()
	if err != nil {
		// Try loading from project root
		err = godotenv.Load("./.env")
		if err != nil {
			fmt.Printf("Warning: .env file not loaded: %v\n", err)
		} else {
			fmt.Println(".env file loaded successfully from ../../.env")
		}
	} else {
		fmt.Println(".env file loaded successfully")
	}

	// Debug: check raw env vars
	fmt.Printf("Raw env - Host: %s, Port: %s, User: %s\n",
		os.Getenv("SYNCOGNIZE_DATABASE_HOST"),
		os.Getenv("SYNCOGNIZE_DATABASE_PORT"),
		os.Getenv("SYNCOGNIZE_DATABASE_USER"))

	v := viper.New()

	// Set up env var binding
	v.SetEnvPrefix("SYNCOGNIZE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind all environment variables explicitly
	bindEnvVars(v)

	// Debug: print what viper sees
	fmt.Printf("Viper sees - Host: %s, Port: %d, User: %s\n",
		v.GetString("database.host"),
		v.GetInt("database.port"),
		v.GetString("database.user"))

	// Optional config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/syncognize/")
	v.AddConfigPath("$HOME/.syncognize")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func bindEnvVars(v *viper.Viper) {
	// Service
	v.BindEnv("service.name")
	v.BindEnv("service.environment")
	v.BindEnv("service.version")

	// gRPC
	v.BindEnv("grpc.port")
	v.BindEnv("grpc.max_recv_msg_size")
	v.BindEnv("grpc.max_send_msg_size")
	v.BindEnv("grpc.connection_timeout")

	// REST
	v.BindEnv("rest.port")
	v.BindEnv("rest.read_timeout")
	v.BindEnv("rest.write_timeout")
	v.BindEnv("rest.cors_origins")

	// MCP
	v.BindEnv("mcp.enabled")
	v.BindEnv("mcp.transport")
	v.BindEnv("mcp.sse_port")

	// Database
	v.BindEnv("database.host")
	v.BindEnv("database.port")
	v.BindEnv("database.user")
	v.BindEnv("database.password")
	v.BindEnv("database.database")
	v.BindEnv("database.ssl_mode")
	v.BindEnv("database.max_open_conns")
	v.BindEnv("database.max_idle_conns")
	v.BindEnv("database.conn_max_lifetime")

	// NATS
	v.BindEnv("nats.url")
	v.BindEnv("nats.enable_jetstream")
	v.BindEnv("nats.max_reconnect")
	v.BindEnv("nats.reconnect_wait")
	v.BindEnv("nats.stream_prefix")
	v.BindEnv("nats.stream_replicas")
	v.BindEnv("nats.max_deliver")
	v.BindEnv("nats.ack_wait")

	// Voyage
	v.BindEnv("voyage.api_key")
	v.BindEnv("voyage.model")
	v.BindEnv("voyage.dimensions")
	v.BindEnv("voyage.batch_size")
	v.BindEnv("voyage.timeout")
	v.BindEnv("voyage.base_url")

	// Gemini
	v.BindEnv("gemini.api_key")
	v.BindEnv("gemini.model")
	v.BindEnv("gemini.max_tokens")
	v.BindEnv("gemini.temperature")
	v.BindEnv("gemini.timeout")

	// Buffer
	v.BindEnv("buffer.time_threshold")
	v.BindEnv("buffer.token_threshold")
	v.BindEnv("buffer.flush_interval")

	// Search
	v.BindEnv("search.default_limit")
	v.BindEnv("search.max_limit")
	v.BindEnv("search.semantic_weight")
	v.BindEnv("search.keyword_weight")
	v.BindEnv("search.min_score")
	v.BindEnv("search.rrf_constant")
	v.BindEnv("search.ppr_damping")
	v.BindEnv("search.ppr_max_iter")

	// Security
	v.BindEnv("security.jwt_secret")
	v.BindEnv("security.jwt_expiry")
	v.BindEnv("security.encryption_key")
	v.BindEnv("security.rate_limit_rps")
	v.BindEnv("security.rate_limit_burst")

	// Logging
	v.BindEnv("logging.level")
	v.BindEnv("logging.development")
	v.BindEnv("logging.encoding")
}

func (c *Config) NewLogger() (*zap.Logger, error) {
	var cfg zap.Config
	if c.Logging.Development {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	if c.Logging.Encoding != "" {
		cfg.Encoding = c.Logging.Encoding
	}

	levelStr := c.Logging.Level

	if levelStr == "" {
		levelStr = "info"
	}

	level, err := zap.ParseAtomicLevel(c.Logging.Level)
	if err != nil {
		return nil, fmt.Errorf("parsing log level: %w", err)
	}
	cfg.Level = level

	return cfg.Build()
}
