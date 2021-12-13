package config

// DatabaseConfig stores database connection options
type DatabaseConfig struct {
	Type         string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	MaxIdleConns int
	MaxOpenConns int
}

// OauthConfig stores oauth service configuration options
type OauthConfig struct {
	AccessTokenLifetime  int
	RefreshTokenLifetime int
	AuthCodeLifetime     int
	InviteCodeLifetime   int
}

// RedisConfig RedisConfig
type RedisConfig struct {
	Addr       string
	Password   string
	DB         int
	Expiration int
}

type AliyunConfig struct {
	RegionId                     string
	AccessKeyId                  string
	AccessKeySecret              string
	OssEndPoint                  string
	OssRoleARN                   string
	OssPermissionDurationSeconds string
	BucketName                   string
	ProjectRootPath              string
	OssAccelerateDomain          string
	OssAccelerateEndPoint        string
}

// Config stores all configuration options
type Config struct {
	Database      DatabaseConfig
	Oauth         OauthConfig
	Redis         RedisConfig
	Aliyun        AliyunConfig
	IsDevelopment bool
	HasDBLog      bool
	LogLocation   string
	SwaggerURL    string
}
