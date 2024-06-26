package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("project name is required")
	}

	if len(os.Args[1]) == 0 {
		log.Fatalln("project name should be non-empty string")
	}

	if len(os.Args) < 3 {
		log.Fatalln("project module name is required")
	}

	if len(os.Args[2]) == 0 {
		log.Fatalln("project name should be non-empty string")
	}

	project := os.Args[1]
	module := os.Args[2]

	dir := filepath.Join(project, "")
	createDir(dir)
	generateGoMod(dir, module)
	generateEnvs(dir)
	generateIgnores(dir)
	generateUtils(dir)
	generateConfig(dir)
	generateRSAKeyPair(dir)
}

func createDir(dir string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("error creating directory: %s", dir)
	}
}

func createFile(file, content string) {
	if err := os.WriteFile(file, []byte(content), os.ModePerm); err != nil {
		log.Fatalf("error creating file: %s", file)
	}
}

func generateRSAKeyPair(dir string) error {
	d := filepath.Join(dir, "keys")
	createDir(d)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	if err := privateKey.Validate(); err != nil {
		return err
	}

	privatePemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	private, err := os.Create(filepath.Join(dir, "keys", "private.pem"))
	if err != nil {
		return err
	}
	defer private.Close()

	if err := pem.Encode(private, privatePemBlock); err != nil {
		return err
	}

	derBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	publicPemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}

	public, err := os.Create(filepath.Join(dir, "keys", "public.pem"))
	if err != nil {
		return err
	}
	defer public.Close()

	if err := pem.Encode(public, publicPemBlock); err != nil {
		return err
	}

	return nil
}

func generateConfig(dir string) {
	env := `package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	// server
	GoMode     string ` + "`" + `mapstructure:"GO_MODE"` + "`" + `
	ServerHost string ` + "`" + `mapstructure:"SERVER_HOST"` + "`" + `
	ServerPort uint16 ` + "`" + `mapstructure:"SERVER_PORT"` + "`" + `
	// database
	DBHost         string ` + "`" + `mapstructure:"DB_HOST"` + "`" + `
	DBName         string ` + "`" + `mapstructure:"DB_NAME"` + "`" + `
	DBPort         uint16 ` + "`" + `mapstructure:"DB_PORT"` + "`" + `
	DBUser         string ` + "`" + `mapstructure:"DB_USER"` + "`" + `
	DBUserPwd      string ` + "`" + `mapstructure:"DB_USER_PWD"` + "`" + `
	DBMinPoolSize  uint16 ` + "`" + `mapstructure:"DB_MIN_POOL_SIZE"` + "`" + `
	DBMaxPoolSize  uint16 ` + "`" + `mapstructure:"DB_MAX_POOL_SIZE"` + "`" + `
	DBQueryTimeout uint16 ` + "`" + `mapstructure:"DB_QUERY_TIMEOUT_SEC"` + "`" + `
	// redis
	RedisHost string ` + "`" + `mapstructure:"REDIS_HOST"` + "`" + `
	RedisPort uint16 ` + "`" + `mapstructure:"REDIS_PORT"` + "`" + `
	RedisPwd  string ` + "`" + `mapstructure:"REDIS_PASSWORD"` + "`" + `
	RedisDB   int    ` + "`" + `mapstructure:"REDIS_DB"` + "`" + `
	// keys
	RSAPrivateKeyPath string ` + "`" + `mapstructure:"RSA_PRIVATE_KEY_PATH"` + "`" + `
	RSAPublicKeyPath  string ` + "`" + `mapstructure:"RSA_PUBLIC_KEY_PATH"` + "`" + `
	// Token
	AccessTokenValiditySec  uint64 ` + "`" + `mapstructure:"ACCESS_TOKEN_VALIDITY_SEC"` + "`" + `
	RefreshTokenValiditySec uint64 ` + "`" + `mapstructure:"REFRESH_TOKEN_VALIDITY_SEC"` + "`" + `
	TokenIssuer             string ` + "`" + `mapstructure:"TOKEN_ISSUER"` + "`" + `
	TokenAudience           string ` + "`" + `mapstructure:"TOKEN_AUDIENCE"` + "`" + `
}

func NewEnv(filename string) *Env {
	env := Env{}
	viper.SetConfigFile(filename)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading environment file", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Error loading environment file", err)
	}

	return &env
}
`
	d := filepath.Join(dir, "config")
	createDir(d)
	createFile(filepath.Join(d, "env.go"), env)
}

func generateUtils(dir string) {
	convertor := `package utils

import (
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

func MapTo[T any, V any](from *V) (*T, error) {
	var to T
	err := copier.Copy(&to, from)
	if err != nil {
		return nil, err
	}
	return &to, nil
}
`
	d := filepath.Join(dir, "utils")
	createDir(d)
	createFile(filepath.Join(d, "convertor.go"), convertor)
}

func generateIgnores(dir string) {
	gitignore := `
 # If you prefer the allow list template instead of the deny list, see community template:
# https://github.com/github/gitignore/blob/main/community/Golang/Go.AllowList.gitignore
#
.DS_Store
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with ` + "`" + `go test -c` + "`" + `
*.test
!Dockerfile.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work
go.work.sum

# Environment varibles
*.env
*.env.test

#keys
keys/*
!keys/*.md
!keys/*.txt
*.pem

__debug*

build
 `
	createFile(filepath.Join(dir, ".gitignore"), gitignore)
}

func generateEnvs(dir string) {
	env := `
	# debug, release, test
GO_MODE=debug

SERVER_HOST=0.0.0.0
SERVER_PORT=8080

DB_HOST=localhost
DB_PORT=27017
DB_NAME=goserver-dev-db
DB_USER=goserver-dev-db-user
DB_USER_PWD=changeit
DB_MIN_POOL_SIZE=2
DB_MAX_POOL_SIZE=5
DB_QUERY_TIMEOUT_SEC=60
DB_ADMIN=admin
DB_ADMIN_PWD=changeit

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=changeit

# 2 DAYS: 172800 Sec
ACCESS_TOKEN_VALIDITY_SEC=172800
# 7 DAYS: 604800 Sec
REFRESH_TOKEN_VALIDITY_SEC=604800
TOKEN_ISSUER=api.goserve.unusualcode.org
TOKEN_AUDIENCE=goserve.unusualcode.org

RSA_PRIVATE_KEY_PATH="keys/private.pem"
RSA_PUBLIC_KEY_PATH="keys/public.pem"
`

	testEnv := `
# debug, release, test
GO_MODE=test

SERVER_HOST=0.0.0.0
SERVER_PORT=8080

DB_HOST=localhost
DB_PORT=27017
DB_NAME=goserver-test-db
DB_USER=goserver-test-db-user
DB_USER_PWD=changeit
DB_MIN_POOL_SIZE=2
DB_MAX_POOL_SIZE=5
DB_QUERY_TIMEOUT_SEC=60
DB_ADMIN=admin
DB_ADMIN_PWD=changeit

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=changeit

# 2 DAYS: 172800 Sec
ACCESS_TOKEN_VALIDITY_SEC=172800
# 7 DAYS: 604800 Sec
REFRESH_TOKEN_VALIDITY_SEC=604800
TOKEN_ISSUER=api.goserve.unusualcode.org
TOKEN_AUDIENCE=goserve.unusualcode.org

RSA_PRIVATE_KEY_PATH="../keys/private.pem"
RSA_PUBLIC_KEY_PATH="../keys/public.pem"
`

	createFile(filepath.Join(dir, ".env"), env)
	createFile(filepath.Join(dir, ".test.env"), testEnv)
}

func generateGoMod(dir, module string) {
	goMod := `module %s

go 1.22.4

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/go-playground/validator/v10 v10.22.0
	github.com/jinzhu/copier v0.4.0
	github.com/spf13/viper v1.19.0
	github.com/unusualcodeorg/goserve v1.0.0
	go.mongodb.org/mongo-driver v1.15.1
)

require (
	github.com/bytedance/sonic v1.11.9 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.4 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/redis/go-redis/v9 v9.5.3 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240424034433-3c2c7870ae76 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`

	goSum := `module %s

go 1.22.4

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/go-playground/validator/v10 v10.22.0
	github.com/jinzhu/copier v0.4.0
	github.com/spf13/viper v1.19.0
	github.com/unusualcodeorg/goserve v1.0.0
	go.mongodb.org/mongo-driver v1.15.1
)

require (
	github.com/bytedance/sonic v1.11.9 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.4 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/redis/go-redis/v9 v9.5.3 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240424034433-3c2c7870ae76 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`
	createFile(filepath.Join(dir, "go.mod"), fmt.Sprintf(goMod, module))
	createFile(filepath.Join(dir, "go.sum"), fmt.Sprintf(goSum, module))
}
