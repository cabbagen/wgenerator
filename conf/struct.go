package conf

type ApplicationConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Cacher   CacherConfig   `yaml:"cacher"`
	Logger   LoggerConfig   `yaml:"logger"`
	Email    EmailConfig    `yaml:"email"`
	Asynq    AsynqConfig    `yaml:"asynq"`
}

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type ServerConfig struct {
	Port             string `yaml:"port"`
	Static           string `yaml:"static"`
	TemplateDir      string `yaml:"templateDir"`
	IsOpenSupportSpa int    `yaml:"isOpenSupportSpa"`
}

type DatabaseConfig struct {
	DBName   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
}

type CacherConfig struct {
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
}

type LoggerConfig struct {
	Output   []string `yaml:"output"`
	Levels   []string `yaml:"levels"`
	File     string   `yaml:"file"`
	Template string   `yaml:"template"`
}

type AsynqConfig struct {
	Redis       AsynqRedisConfig `yaml:"redis"`
	Concurrency int              `yaml:"concurrency"`
}

type AsynqRedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
