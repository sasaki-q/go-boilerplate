package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	sm "github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	sm.NewSecret(stack, jsii.String("secret-manager"), &sm.SecretProps{
		SecretName: jsii.String("boilerplate-db"),
		SecretObjectValue: &map[string]awscdk.SecretValue{
			"host":     awscdk.SecretValue_UnsafePlainText(jsii.String("boilerplate_db")),
			"dbname":   awscdk.SecretValue_UnsafePlainText(jsii.String("boilerplate_db")),
			"user":     awscdk.SecretValue_UnsafePlainText(jsii.String("postgres")),
			"password": awscdk.SecretValue_UnsafePlainText(jsii.String("postgres")),
			"port":     awscdk.SecretValue_UnsafePlainText(jsii.String("5432")),
		},
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	var (
		abn     = app.Node().TryGetContext(jsii.String("ABN"))
		env     = app.Node().TryGetContext(jsii.String("ENV"))
		project = app.Node().TryGetContext(jsii.String("P"))
	)

	if abn == nil || env == nil || project == nil {
		panic("please pass context")
	}

	var (
		a = fmt.Sprintf("%s", abn)
		e = fmt.Sprintf("%s", env)
		p = fmt.Sprintf("%s", project)
	)
	print(e)
	print(p)

	awscdk.Tags_Of(app).Add(jsii.String("Env"), &e, nil)
	awscdk.Tags_Of(app).Add(jsii.String("Project"), &p, nil)
	NewCdkStack(app, fmt.Sprintf("%sStack", p),
		&CdkStackProps{
			awscdk.StackProps{
				Env: myenv(),
				Synthesizer: awscdk.NewDefaultStackSynthesizer(
					&awscdk.DefaultStackSynthesizerProps{FileAssetsBucketName: &a},
				),
			},
		},
	)

	app.Synth(nil)
}

func myenv() *awscdk.Environment { return nil }
