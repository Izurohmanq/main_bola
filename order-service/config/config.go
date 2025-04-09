package config

import (
	"order-service/common/utils"
	"os"

	"github.com/sirupsen/logrus"
	_ "github.com/spf13/viper/remote"
)

var Config AppConfig

// ini sesuai dengan config.json yang dibuat
type AppConfig struct {
	Port                       int             `json:"port"`
	AppName                    string          `json:"appName"`
	AppEnv                     string          `json:"appEnv"`
	SignatureKey               string          `json:"signatureKey"`
	Database                   Database        `json:"database"`
	RateLimiterMaxRequest      float64         `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond      int             `json:"rateLimiterTimeSecond"`
	JwtSecretKey               string          `json:"jwtSecretKey"`
	JwtExpirationTime          int             `json:"jwtExpirationTime"`
	InternalService            InternalService `json:"internalService"`
	GCSType                    string          `json:"gcsType"`
	GCSProjectID               string          `json:"gcsProjectID"`
	GCSPrivateKeyID            string          `json:"gcsPrivateKeyID"`
	GCSPrivateKey              string          `json:"gcsPrivateKey"`
	GCSClientEmail             string          `json:"gcsClientEmail"`
	GCSClientID                string          `json:"gcsClientID"`
	GCSAuthURI                 string          `json:"gcsAuthURI"`
	GCSTokenURI                string          `json:"gcsTokenURI"`
	GCSAuthProviderX509CertURL string          `json:"gcsAuthProviderX509CertURL"`
	GCSClientX509CertURL       string          `json:"gcsClientX509CertURL"`
	GCSUniverseDomain          string          `json:"gcsUniverseDomain"`
	GCSBucketName              string          `json:"gcsBucketName"`
	Kafka                      Kafka           `json:"kafka"`
}

type Database struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	Name                  string `json:"name"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	MaxOpenConnection     int    `json:"maxOpenConnection"`
	MaxLifetimeConnection int    `json:"maxLifetimeConnection"`
	MaxIdleConnection     int    `json:"maxIdleConnection"`
	MaxIdleTime           int    `json:"maxIdleTime"`
}

type InternalService struct {
	User    User    `json:"user"`
	Field   Field   `json:"field"`
	Payment Payment `json:"payment"`
}

type Kafka struct {
	Brokers               []string `json:"brokers"`
	TimeoutInMs           int      `json:"timeoutInMs"`
	MaxRetry              int      `json:"maxRetry"`
	MaxWaitTimeInMs       int      `json:"maxWaitTimeInMs"`
	MaxProcessingTimeInMs int      `json:"maxProcessingTimeInMs"`
	BackOffTimeInMs       int      `json:"backoffTimeInMs"`
	Topics                []string `json:"topics"`
	GroupID               string   `json:"groupID"`
}

type User struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

type Field struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

type Payment struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

func Init() {
	err := utils.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("failed to bind config: %v", err)
		err = utils.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
