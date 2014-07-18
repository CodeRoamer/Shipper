package setting

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/Unknwon/com"
	"github.com/gogits/cache"
	"github.com/gogits/session"

	"github.com/coderoamer/shipper/modules/log"
)

type Scheme string

const (
	HTTP Scheme = "http"
	HTTPS Scheme = "https"
)

var (
	// App settings.
	AppVer  string // not in ini
	AppName string
	AppLogo string
	AppUrl  string // not in ini

	// Server settings.
	Protocol           Scheme
	Domain             string
	HttpAddr, HttpPort string
	SshPort            int
	OfflineMode        bool
	DisableRouterLog   bool
	CertFile, KeyFile  string
	StaticRootPath     string

	// Security settings.
	InstallLock          bool
	SecretKey            string
	LogInRememberDays    int
	CookieUserName       string
	CookieRememberName   string
	ReverseProxyAuthUser string

	// Webhook settings.
	WebhookTaskInterval   int
	WebhookDeliverTimeout int

	// Log settings.
	LogRootPath string
	LogModes    []string
	LogConfigs  []string

	// Cache settings.
	Cache        cache.Cache
	CacheAdapter string
	CacheConfig  string

	EnableRedis    bool
	EnableMemcache bool

	// Session settings.
	SessionProvider string
	SessionConfig   *session.Config
	SessionManager  *session.Manager

	// Global setting objects.
	Cfg        *goconfig.ConfigFile
	ProdMode   bool
	RunUser    string
)

// Load a file according to the given string
func ExecPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return p, nil
}

// WorkDir returns absolute path of work directory.
func WorkDir() (string, error) {
	execPath, err := ExecPath()
	return path.Dir(strings.Replace(execPath, "\\", "/", -1)), err
}

// NewConfigContext initializes configuration context.
// NOTE: do not print any log except error.
func NewConfigContext() {
	workDir, err := WorkDir()
	if err != nil {
		log.Fatal("Fail to get work directory: %v", err)
	}

	cfgPath := path.Join(workDir, "conf/app.ini")
	if com.IsFile(cfgPath) {
		Cfg, err = goconfig.LoadConfigFile(cfgPath)
		if err != nil {
			log.Fatal("Fail to load 'conf/app.ini': %v", err)
		}
	} else {
		log.Warn("No 'conf/app.ini' found")
	}

	AppName = Cfg.MustValue("", "APP_NAME", "Shipper: Docker Web UI written in GO")
	AppLogo = Cfg.MustValue("", "APP_LOGO", "img/favicon.png")
	AppUrl = Cfg.MustValue("server", "ROOT_URL", "http://localhost:1212")

	Protocol = HTTP
	if Cfg.MustValue("server", "PROTOCOL") == "https" {
		Protocol = HTTPS
		CertFile = Cfg.MustValue("server", "CERT_FILE")
		KeyFile = Cfg.MustValue("server", "KEY_FILE")
	}

	Domain = Cfg.MustValue("server", "DOMAIN", "localhost")
	HttpAddr = Cfg.MustValue("server", "HTTP_ADDR", "0.0.0.0")
	HttpPort = Cfg.MustValue("server", "HTTP_PORT", "1212")
	SshPort = Cfg.MustInt("server", "SSH_PORT", 22)
	OfflineMode = Cfg.MustBool("server", "OFFLINE_MODE")
	DisableRouterLog = Cfg.MustBool("server", "DISABLE_ROUTER_LOG")
	StaticRootPath = Cfg.MustValue("server", "STATIC_ROOT_PATH", workDir)
	LogRootPath = Cfg.MustValue("log", "ROOT_PATH", path.Join(workDir, "log"))

	InstallLock = Cfg.MustBool("security", "INSTALL_LOCK")
	SecretKey = Cfg.MustValue("security", "SECRET_KEY")
	LogInRememberDays = Cfg.MustInt("security", "LOGIN_REMEMBER_DAYS")
	CookieUserName = Cfg.MustValue("security", "COOKIE_USERNAME")
	CookieRememberName = Cfg.MustValue("security", "COOKIE_REMEMBER_NAME")
	ReverseProxyAuthUser = Cfg.MustValue("security", "REVERSE_PROXY_AUTHENTICATION_USER", "X-WEBAUTH-USER")

	RunUser = Cfg.MustValue("", "RUN_USER")
	curUser := os.Getenv("USER")
	if len(curUser) == 0 {
		curUser = os.Getenv("USERNAME")
	}
	// Does not check run user when the install lock is off.
	if InstallLock && RunUser != curUser {
		log.Fatal("Expect user(%s) but current user is: %s", RunUser, curUser)
	}

}

var Service struct {
		RegisterEmailConfirm   bool
		DisableRegistration    bool
		RequireSignInView      bool
		EnableCacheAvatar      bool
		EnableNotifyMail       bool
		EnableReverseProxyAuth bool
		ActiveCodeLives        int
		ResetPwdCodeLives      int
}


func newService() {
	Service.ActiveCodeLives = Cfg.MustInt("service", "ACTIVE_CODE_LIVE_MINUTES", 180)
	Service.ResetPwdCodeLives = Cfg.MustInt("service", "RESET_PASSWD_CODE_LIVE_MINUTES", 180)
	Service.DisableRegistration = Cfg.MustBool("service", "DISABLE_REGISTRATION")
	Service.RequireSignInView = Cfg.MustBool("service", "REQUIRE_SIGNIN_VIEW")
	Service.EnableCacheAvatar = Cfg.MustBool("service", "ENABLE_CACHE_AVATAR")
	Service.EnableReverseProxyAuth = Cfg.MustBool("service", "ENABLE_REVERSE_PROXY_AUTHENTICATION")
}

var logLevels = map[string]string{
	"Trace":    "0",
	"Debug":    "1",
	"Info":     "2",
	"Warn":     "3",
	"Error":    "4",
	"Critical": "5",
}


func newLogService() {
	log.Info("%s %s", AppName, AppVer)

	// Get and check log mode.
	LogModes = strings.Split(Cfg.MustValue("log", "MODE", "console"), ",")
	LogConfigs = make([]string, len(LogModes))
	for i, mode := range LogModes {
		mode = strings.TrimSpace(mode)
		modeSec := "log." + mode
		if _, err := Cfg.GetSection(modeSec); err != nil {
			log.Fatal("Unknown log mode: %s", mode)
		}

		// Log level.
		levelName := Cfg.MustValueRange("log."+mode, "LEVEL", "Trace",
			[]string{"Trace", "Debug", "Info", "Warn", "Error", "Critical"})
		level, ok := logLevels[levelName]
		if !ok {
			log.Fatal("Unknown log level: %s", levelName)
		}

		// Generate log configuration.
		switch mode {
		case "console":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s}`, level)
		case "file":
			logPath := Cfg.MustValue(modeSec, "FILE_NAME", path.Join(LogRootPath, "gogs.log"))
			os.MkdirAll(path.Dir(logPath), os.ModePerm)
			LogConfigs[i] = fmt.Sprintf(
				`{"level":%s,"filename":"%s","rotate":%v,"maxlines":%d,"maxsize":%d,"daily":%v,"maxdays":%d}`, level,
				logPath,
				Cfg.MustBool(modeSec, "LOG_ROTATE", true),
				Cfg.MustInt(modeSec, "MAX_LINES", 1000000),
					1<<uint(Cfg.MustInt(modeSec, "MAX_SIZE_SHIFT", 28)),
				Cfg.MustBool(modeSec, "DAILY_ROTATE", true),
				Cfg.MustInt(modeSec, "MAX_DAYS", 7))
		case "database":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s,"driver":"%s","conn":"%s"}`, level,
				Cfg.MustValue(modeSec, "DRIVER"),
				Cfg.MustValue(modeSec, "CONN"))
		}

		log.NewLogger(Cfg.MustInt64("log", "BUFFER_LEN", 10000), mode, LogConfigs[i])
		log.Info("Log Mode: %s(%s)", strings.Title(mode), levelName)
	}
}

func newCacheService() {

}

func newSessionService() {

}

func newMailService() {

}

func newRegisterMailService() {

}

func newNotifyMailService() {

}

func newWebhookService() {

}

func NewServices() {
	newService()
	newLogService()
	newCacheService()
	newSessionService()
	newMailService()
	newRegisterMailService()
	newNotifyMailService()
	newWebhookService()
}
