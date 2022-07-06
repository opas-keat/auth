package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper

	errorHandler fiber.ErrorHandler
	fiber        *fiber.Config
}


var (
	KEY = "4244467890218023942835864651981516513207895156132164643213169699"
	// passphrase    = "AAACCCDDDYYUURRS"
)

var defaultErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Set error message
	message := err.Error()

	// Check if it's a fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// TODO: Check return type for the client, JSON, HTML, YAML or any other (API vs web)

	// Return HTTP response
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	c.Status(code)

	// Render default error view
	err = c.Render("errors/"+strconv.Itoa(code), fiber.Map{"message": message})
	if err != nil {
		return c.SendString(message)
	}
	return err
}

func New() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// Select the .env file
	config.SetConfigName(".env")
	config.SetConfigType("dotenv")
	config.AddConfigPath(".")

	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}

	config.SetErrorHandler(defaultErrorHandler)

	// TODO: Logger (Maybe a different zap object)

	// TODO: Add APP_KEY generation

	// TODO: Write changes to configuration file

	// Set Fiber configurations
	config.setFiberConfig()

	return config
}

func (config *Config) SetErrorHandler(errorHandler fiber.ErrorHandler) {
	config.errorHandler = errorHandler
}

func (config *Config) setDefaults() {
	// Set default App configuration
	config.SetDefault("APP_ADDR", "127.0.0.1:8080")
	config.SetDefault("APP_ENV", "local")

	// Set default database configuration
	config.SetDefault("DB_DRIVER", "mysql")
	config.SetDefault("DB_HOST", "localhost")
	config.SetDefault("DB_USERNAME", "fiber")
	config.SetDefault("DB_PASSWORD", "password")
	config.SetDefault("DB_PORT", 3306)
	config.SetDefault("DB_DATABASE", "boilerplate")

	// Set default hasher configuration
	// config.SetDefault("HASHER_DRIVER", "argon2id")
	// config.SetDefault("HASHER_MEMORY", 131072)
	// config.SetDefault("HASHER_ITERATIONS", 4)
	// config.SetDefault("HASHER_PARALLELISM", 4)
	// config.SetDefault("HASHER_SALTLENGTH", 16)
	// config.SetDefault("HASHER_KEYLENGTH", 32)
	// config.SetDefault("HASHER_ROUNDS", bcrypt.DefaultRounds)

	// Set default session configuration
	config.SetDefault("SESSION_PROVIDER", "mysql")
	config.SetDefault("SESSION_KEYPREFIX", "session")
	config.SetDefault("SESSION_HOST", "localhost")
	config.SetDefault("SESSION_PORT", 3306)
	config.SetDefault("SESSION_USERNAME", "fiber")
	config.SetDefault("SESSION_PASSWORD", "secret")
	config.SetDefault("SESSION_DATABASE", "boilerplate")
	config.SetDefault("SESSION_TABLENAME", "sessions")
	config.SetDefault("SESSION_LOOKUP", "cookie:session_id")
	config.SetDefault("SESSION_DOMAIN", "")
	config.SetDefault("SESSION_SAMESITE", "Lax")
	config.SetDefault("SESSION_EXPIRATION", "12h")
	config.SetDefault("SESSION_SECURE", false)
	config.SetDefault("SESSION_GCINTERVAL", "1m")

	// Set default Fiber configuration
	config.SetDefault("FIBER_PREFORK", false)
	config.SetDefault("FIBER_SERVERHEADER", "")
	config.SetDefault("FIBER_STRICTROUTING", false)
	config.SetDefault("FIBER_CASESENSITIVE", false)
	config.SetDefault("FIBER_IMMUTABLE", false)
	config.SetDefault("FIBER_UNESCAPEPATH", false)
	config.SetDefault("FIBER_ETAG", false)
	config.SetDefault("FIBER_BODYLIMIT", 4194304)
	config.SetDefault("FIBER_CONCURRENCY", 262144)
	config.SetDefault("FIBER_VIEWS", "html")
	config.SetDefault("FIBER_VIEWS_DIRECTORY", "resources/views")
	config.SetDefault("FIBER_VIEWS_RELOAD", false)
	config.SetDefault("FIBER_VIEWS_DEBUG", false)
	config.SetDefault("FIBER_VIEWS_LAYOUT", "embed")
	config.SetDefault("FIBER_VIEWS_DELIMS_L", "{{")
	config.SetDefault("FIBER_VIEWS_DELIMS_R", "}}")
	config.SetDefault("FIBER_READTIMEOUT", 0)
	config.SetDefault("FIBER_WRITETIMEOUT", 0)
	config.SetDefault("FIBER_IDLETIMEOUT", 0)
	config.SetDefault("FIBER_READBUFFERSIZE", 4096)
	config.SetDefault("FIBER_WRITEBUFFERSIZE", 4096)
	config.SetDefault("FIBER_COMPRESSEDFILESUFFIX", ".fiber.gz")
	config.SetDefault("FIBER_PROXYHEADER", "")
	config.SetDefault("FIBER_GETONLY", false)
	config.SetDefault("FIBER_DISABLEKEEPALIVE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTDATE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTCONTENTTYPE", false)
	config.SetDefault("FIBER_DISABLEHEADERNORMALIZING", false)
	config.SetDefault("FIBER_DISABLESTARTUPMESSAGE", false)
	config.SetDefault("FIBER_REDUCEMEMORYUSAGE", false)

	// Set default Custom Access Logger middleware configuration
	config.SetDefault("MW_ACCESS_LOGGER_ENABLED", true)
	config.SetDefault("MW_ACCESS_LOGGER_TYPE", "console")
	config.SetDefault("MW_ACCESS_LOGGER_FILENAME", "access.log")
	config.SetDefault("MW_ACCESS_LOGGER_MAXSIZE", 500)
	config.SetDefault("MW_ACCESS_LOGGER_MAXAGE", 28)
	config.SetDefault("MW_ACCESS_LOGGER_MAXBACKUPS", 3)
	config.SetDefault("MW_ACCESS_LOGGER_LOCALTIME", false)
	config.SetDefault("MW_ACCESS_LOGGER_COMPRESS", false)

	// Set default Force HTTPS middleware configuration
	config.SetDefault("MW_FORCE_HTTPS_ENABLED", false)

	// Set default Force trailing slash middleware configuration
	config.SetDefault("MW_FORCE_TRAILING_SLASH_ENABLED", false)

	// Set default HSTS middleware configuration
	config.SetDefault("MW_HSTS_ENABLED", false)
	config.SetDefault("MW_HSTS_MAXAGE", 31536000)
	config.SetDefault("MW_HSTS_INCLUDESUBDOMAINS", true)
	config.SetDefault("MW_HSTS_PRELOAD", false)

	// Set default Suppress WWW middleware configuration
	config.SetDefault("MW_SUPPRESS_WWW_ENABLED", true)

	// Set default Fiber Cache middleware configuration
	config.SetDefault("MW_FIBER_CACHE_ENABLED", false)
	config.SetDefault("MW_FIBER_CACHE_EXPIRATION", "1m")
	config.SetDefault("MW_FIBER_CACHE_CACHECONTROL", false)

	// Set default Fiber Compress middleware configuration
	config.SetDefault("MW_FIBER_COMPRESS_ENABLED", false)
	config.SetDefault("MW_FIBER_COMPRESS_LEVEL", 0)

	// Set default Fiber CORS middleware configuration
	config.SetDefault("MW_FIBER_CORS_ENABLED", false)
	config.SetDefault("MW_FIBER_CORS_ALLOWORIGINS", "*")
	config.SetDefault("MW_FIBER_CORS_ALLOWMETHODS", "GET,POST,HEAD,PUT,DELETE,PATCH")
	config.SetDefault("MW_FIBER_CORS_ALLOWHEADERS", "")
	config.SetDefault("MW_FIBER_CORS_ALLOWCREDENTIALS", false)
	config.SetDefault("MW_FIBER_CORS_EXPOSEHEADERS", "")
	config.SetDefault("MW_FIBER_CORS_MAXAGE", 0)

	// Set default Fiber CSRF middleware configuration
	config.SetDefault("MW_FIBER_CSRF_ENABLED", false)
	config.SetDefault("MW_FIBER_CSRF_TOKENLOOKUP", "header:X-CSRF-Token")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_NAME", "_csrf")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_SAMESITE", "Strict")
	config.SetDefault("MW_FIBER_CSRF_COOKIE_EXPIRES", "24h")
	config.SetDefault("MW_FIBER_CSRF_CONTEXTKEY", "csrf")

	// Set default Fiber ETag middleware configuration
	config.SetDefault("MW_FIBER_ETAG_ENABLED", false)
	config.SetDefault("MW_FIBER_ETAG_WEAK", false)

	// Set default Fiber Expvar middleware configuration
	config.SetDefault("MW_FIBER_EXPVAR_ENABLED", false)

	// Set default Fiber Favicon middleware configuration
	config.SetDefault("MW_FIBER_FAVICON_ENABLED", false)
	config.SetDefault("MW_FIBER_FAVICON_FILE", "")
	config.SetDefault("MW_FIBER_FAVICON_CACHECONTROL", "public, max-age=31536000")

	// Set default Fiber Limiter middleware configuration
	config.SetDefault("MW_FIBER_LIMITER_ENABLED", false)
	config.SetDefault("MW_FIBER_LIMITER_MAX", 5)
	config.SetDefault("MW_FIBER_LIMITER_DURATION", "1m")

	// Set default Fiber Monitor middleware configuration
	config.SetDefault("MW_FIBER_MONITOR_ENABLED", false)

	// Set default Fiber Pprof middleware configuration
	config.SetDefault("MW_FIBER_PPROF_ENABLED", false)

	// Set default Fiber Recover middleware configuration
	config.SetDefault("MW_FIBER_RECOVER_ENABLED", true)

	// Set default Fiber RequestID middleware configuration
	config.SetDefault("MW_FIBER_REQUESTID_ENABLED", false)
	config.SetDefault("MW_FIBER_REQUESTID_HEADER", "X-Request-ID")
	config.SetDefault("MW_FIBER_REQUESTID_CONTEXTKEY", "requestid")
}

func (config *Config) setFiberConfig() {
	config.fiber = &fiber.Config{
		Prefork:       config.GetBool("FIBER_PREFORK"),
		ServerHeader:  config.GetString("FIBER_SERVERHEADER"),
		StrictRouting: config.GetBool("FIBER_STRICTROUTING"),
		CaseSensitive: config.GetBool("FIBER_CASESENSITIVE"),
		Immutable:     config.GetBool("FIBER_IMMUTABLE"),
		UnescapePath:  config.GetBool("FIBER_UNESCAPEPATH"),
		ETag:          config.GetBool("FIBER_ETAG"),
		BodyLimit:     config.GetInt("FIBER_BODYLIMIT"),
		Concurrency:   config.GetInt("FIBER_CONCURRENCY"),
		// Views:                     config.getFiberViewsEngine(),
		ReadTimeout:               config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:              config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:               config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:            config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:           config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix:      config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:               config.GetString("FIBER_PROXYHEADER"),
		GETOnly:                   config.GetBool("FIBER_GETONLY"),
		ErrorHandler:              config.errorHandler,
		DisableKeepalive:          config.GetBool("FIBER_DISABLEKEEPALIVE"),
		DisableDefaultDate:        config.GetBool("FIBER_DISABLEDEFAULTDATE"),
		DisableDefaultContentType: config.GetBool("FIBER_DISABLEDEFAULTCONTENTTYPE"),
		DisableHeaderNormalizing:  config.GetBool("FIBER_DISABLEHEADERNORMALIZING"),
		DisableStartupMessage:     config.GetBool("FIBER_DISABLESTARTUPMESSAGE"),
		ReduceMemoryUsage:         config.GetBool("FIBER_REDUCEMEMORYUSAGE"),
		JSONEncoder:               json.Marshal,
		JSONDecoder:               json.Unmarshal,
	}
}

func (config *Config) GetFiberConfig() *fiber.Config {
	return config.fiber
}

// func (config *Config) GetSessionConfig() session.Config {
// 	var provider fsession.Provider
// 	switch strings.ToLower(config.GetString("SESSION_PROVIDER")) {
// 	case "redis":
// 		sessionProvider, err := redis.New(redis.Config{
// 			KeyPrefix: config.GetString("SESSION_KEYPREFIX"),
// 			Addr:      config.GetString("SESSION_HOST") + ":" + config.GetString("SESSION_PORT"),
// 			Password:  config.GetString("SESSION_PASSWORD"),
// 			DB:        config.GetInt("SESSION_DATABASE"),
// 		})
// 		if err != nil {
// 			fmt.Println("failed to initialized redis session provider:", err.Error())
// 			break
// 		}
// 		provider = sessionProvider
// 		break
// 	}

// 	return session.Config{
// 		Lookup:     config.GetString("SESSION_LOOKUP"),
// 		Secure:     config.GetBool("SESSION_SECURE"),
// 		Domain:     config.GetString("SESSION_DOMAIN"),
// 		SameSite:   config.GetString("SESSION_SAMESITE"),
// 		Expiration: config.GetDuration("SESSION_EXPIRATION"),
// 		Provider:   provider,
// 		GCInterval: config.GetDuration("SESSION_GCINTERVAL"),
// 	}
// }
