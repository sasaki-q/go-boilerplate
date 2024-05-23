package services

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"go.uber.org/zap"
)

type AWSService struct {
	l        *zap.Logger
	smClient *sm.Client
}

type DBValue struct {
	Password string `json:"password" validate:"required"`
	DBName   string `json:"dbname" validate:"required"`
	Host     string `json:"host" validate:"required"`
	User     string `json:"user" validate:"required"`
}

func NewAWSService(ctx context.Context, l *zap.Logger) *AWSService {
	cfg, err := newConfig(ctx)
	if err != nil {
		l.Sugar().Errorf("err cannot load config. / %v", err)
		panic(fmt.Sprintf("err cannot load config. / %v", err))
	}

	return &AWSService{
		l:        l,
		smClient: sm.NewFromConfig(cfg),
	}
}

func (s *AWSService) GetSecretValue(ctx context.Context, ID string) (*sm.GetSecretValueOutput, error) {
	return s.smClient.GetSecretValue(ctx, &sm.GetSecretValueInput{SecretId: &ID})
}

func newConfig(ctx context.Context) (aws.Config, error) {
	if os.Getenv("ENV") == "dev" {
		p := credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")

		return config.LoadDefaultConfig(
			ctx,
			config.WithRegion("ap-northeast-1"),
			config.WithCredentialsProvider(p),
		)
	}

	return config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
}
