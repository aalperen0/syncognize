package config

import (
	"fmt"
	"os"
	"strconv"
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

// DSN returns postgres connection string
// func (d DatabaseConfig) DSN() string {
// 	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
// 		d.User, d.Password, d.Host, d.Port, d.Database, d.SSLMode)
// }

// Load reads config from env vars (SYNCOGNIZE_ prefix) and optional config file
func Load() (*Config, error) {
	v := viper.New()

	setDefaults(v)

	v.SetEnvPrefix("SYNCOGNIZE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

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

	level, err := zap.ParseAtomicLevel(c.Logging.Level)
	if err != nil {
		return nil, fmt.Errorf("parsing log level: %w", err)
	}
	cfg.Level = level

	return cfg.Build()
}

func setDefaults(v *viper.Viper) {

	if err := godotenv.Load(); err != nil {
		zap.L().Info("No .env files found, using environment variables")
	}

	getEnv := func(key string) string {
		return os.Getenv(key)
	}

	getEnvInt := func(key string) int {
		if val := os.Getenv(key); val != "" {
			if i, err := strconv.Atoi(val); err == nil {
				return i
			}
		}
		return 0
	}

	getEnvBool := func(key string) bool {
		if val := os.Getenv(key); val != "" {
			if b, err := strconv.ParseBool(val); err == nil {
				return b
			}
		}
		return false
	}

	getEnvFloat := func(key string) float32 {
		if val := os.Getenv(key); val != "" {
			if f, err := strconv.ParseFloat(val, 32); err == nil {
				return float32(f)
			}
		}
		return 0
	}

	// Service
	v.SetDefault("service.name", getEnv("SYNCOGNIZE_SERVICE_NAME"))
	v.SetDefault("service.environment", getEnv("SYNCOGNIZE_SERVICE_ENVIRONMENT"))
	v.SetDefault("service.version", getEnv("SYNCOGNIZE_SERVICE_VERSION"))

	// gRPC
	v.SetDefault("grpc.port", getEnvInt("SYNCOGNIZE_GRPC_PORT"))
	v.SetDefault("grpc.max_recv_msg_size", getEnvInt("SYNCOGNIZE_GRPC_MAX_RECV_MSG_SIZE"))
	v.SetDefault("grpc.max_send_msg_size", getEnvInt("SYNCOGNIZE_GRPC_MAX_SEND_MSG_SIZE"))
	v.SetDefault("grpc.connection_timeout", getEnv("SYNCOGNIZE_GRPC_CONNECTION_TIMEOUT"))

	// REST
	v.SetDefault("rest.port", getEnvInt("SYNCOGNIZE_REST_PORT"))
	v.SetDefault("rest.read_timeout", getEnv("SYNCOGNIZE_REST_READ_TIMEOUT"))
	v.SetDefault("rest.write_timeout", getEnv("SYNCOGNIZE_REST_WRITE_TIMEOUT"))
	v.SetDefault("rest.cors_origins", strings.Split(getEnv("SYNCOGNIZE_REST_CORS_ORIGINS"), ","))

	// MCP
	v.SetDefault("mcp.enabled", getEnvBool("SYNCOGNIZE_MCP_ENABLED"))
	v.SetDefault("mcp.transport", getEnv("SYNCOGNIZE_MCP_TRANSPORT"))
	v.SetDefault("mcp.sse_port", getEnvInt("SYNCOGNIZE_MCP_SSE_PORT"))

	// Database
	v.SetDefault("database.host", getEnv("SYNCOGNIZE_DATABASE_HOST"))
	v.SetDefault("database.port", getEnvInt("SYNCOGNIZE_DATABASE_PORT"))
	v.SetDefault("database.user", getEnv("SYNCOGNIZE_DATABASE_USER"))
	v.SetDefault("database.password", getEnv("SYNCOGNIZE_DATABASE_PASSWORD"))
	v.SetDefault("database.database", getEnv("SYNCOGNIZE_DATABASE_NAME"))
	v.SetDefault("database.ssl_mode", getEnv("SYNCOGNIZE_DATABASE_SSL_MODE"))
	v.SetDefault("database.max_open_conns", getEnvInt("SYNCOGNIZE_DATABASE_MAX_OPEN_CONNS"))
	v.SetDefault("database.max_idle_conns", getEnvInt("SYNCOGNIZE_DATABASE_MAX_IDLE_CONNS"))
	v.SetDefault("database.conn_max_lifetime", getEnv("SYNCOGNIZE_DATABASE_CONN_MAX_LIFETIME"))

	// NATS
	v.SetDefault("nats.url", getEnv("SYNCOGNIZE_NATS_URL"))
	v.SetDefault("nats.enable_jetstream", getEnvBool("SYNCOGNIZE_NATS_ENABLE_JETSTREAM"))
	v.SetDefault("nats.max_reconnect", getEnvInt("SYNCOGNIZE_NATS_MAX_RECONNECT"))
	v.SetDefault("nats.reconnect_wait", getEnv("SYNCOGNIZE_NATS_RECONNECT_WAIT"))
	v.SetDefault("nats.stream_prefix", getEnv("SYNCOGNIZE_NATS_STREAM_PREFIX"))
	v.SetDefault("nats.stream_replicas", getEnvInt("SYNCOGNIZE_NATS_STREAM_REPLICAS"))
	v.SetDefault("nats.max_deliver", getEnvInt("SYNCOGNIZE_NATS_MAX_DELIVER"))
	v.SetDefault("nats.ack_wait", getEnv("SYNCOGNIZE_NATS_ACK_WAIT"))

	// Voyage AI
	v.SetDefault("voyage.api_key", getEnv("SYNCOGNIZE_VOYAGE_API_KEY"))
	v.SetDefault("voyage.model", getEnv("SYNCOGNIZE_VOYAGE_MODEL"))
	v.SetDefault("voyage.dimensions", getEnvInt("SYNCOGNIZE_VOYAGE_DIMENSIONS"))
	v.SetDefault("voyage.batch_size", getEnvInt("SYNCOGNIZE_VOYAGE_BATCH_SIZE"))
	v.SetDefault("voyage.timeout", getEnv("SYNCOGNIZE_VOYAGE_TIMEOUT"))
	v.SetDefault("voyage.base_url", getEnv("SYNCOGNIZE_VOYAGE_BASE_URL"))

	// Gemini
	v.SetDefault("gemini.api_key", getEnv("SYNCOGNIZE_GEMINI_API_KEY"))
	v.SetDefault("gemini.model", getEnv("SYNCOGNIZE_GEMINI_MODEL"))
	v.SetDefault("gemini.max_tokens", getEnvInt("SYNCOGNIZE_GEMINI_MAX_TOKENS"))
	v.SetDefault("gemini.temperature", getEnvFloat("SYNCOGNIZE_GEMINI_TEMPERATURE"))
	v.SetDefault("gemini.timeout", getEnv("SYNCOGNIZE_GEMINI_TIMEOUT"))

	// Buffer
	v.SetDefault("buffer.time_threshold", getEnv("SYNCOGNIZE_BUFFER_TIME_THRESHOLD"))
	v.SetDefault("buffer.token_threshold", getEnvInt("SYNCOGNIZE_BUFFER_TOKEN_THRESHOLD"))
	v.SetDefault("buffer.flush_interval", getEnv("SYNCOGNIZE_BUFFER_FLUSH_INTERVAL"))

	// Search
	v.SetDefault("search.default_limit", getEnvInt("SYNCOGNIZE_SEARCH_DEFAULT_LIMIT"))
	v.SetDefault("search.max_limit", getEnvInt("SYNCOGNIZE_SEARCH_MAX_LIMIT"))
	v.SetDefault("search.semantic_weight", getEnvFloat("SYNCOGNIZE_SEARCH_SEMANTIC_WEIGHT"))
	v.SetDefault("search.keyword_weight", getEnvFloat("SYNCOGNIZE_SEARCH_KEYWORD_WEIGHT"))
	v.SetDefault("search.min_score", getEnvFloat("SYNCOGNIZE_SEARCH_MIN_SCORE"))
	v.SetDefault("search.rrf_constant", getEnvInt("SYNCOGNIZE_SEARCH_RRF_CONSTANT"))
	v.SetDefault("search.ppr_damping", getEnvFloat("SYNCOGNIZE_SEARCH_PPR_DAMPING"))
	v.SetDefault("search.ppr_max_iter", getEnvInt("SYNCOGNIZE_SEARCH_PPR_MAX_ITER"))

	// Security
	v.SetDefault("security.jwt_secret", getEnv("SYNCOGNIZE_SECURITY_JWT_SECRET"))
	v.SetDefault("security.jwt_expiry", getEnv("SYNCOGNIZE_SECURITY_JWT_EXPIRY"))
	v.SetDefault("security.encryption_key", getEnv("SYNCOGNIZE_SECURITY_ENCRYPTION_KEY"))
	v.SetDefault("security.rate_limit_rps", getEnvInt("SYNCOGNIZE_SECURITY_RATE_LIMIT_RPS"))
	v.SetDefault("security.rate_limit_burst", getEnvInt("SYNCOGNIZE_SECURITY_RATE_LIMIT_BURST"))

	// Logging
	v.SetDefault("logging.level", getEnv("SYNCOGNIZE_LOGGING_LEVEL"))
	v.SetDefault("logging.development", getEnvBool("SYNCOGNIZE_LOGGING_DEVELOPMENT"))
	v.SetDefault("logging.encoding", getEnv("SYNCOGNIZE_LOGGING_ENCODING"))
}
