package resolver

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	getter "github.com/hashicorp/go-getter"
	urlhelper "github.com/hashicorp/go-getter/helper/url"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	aferoS3 "github.com/wreulicke/afero-s3"
)

var (
	s3Detector = new(getter.S3Detector)
)

// S3Resolver resolves s3 URLs e.g. s3.amazonaws.com/replicated-analyze-testing/ABCD1234/bundle.tgz
type S3Resolver struct {
}

func (g *S3Resolver) Detect(name string) (string, bool, error) {
	forced, ok, err := s3Detector.Detect(name, "")
	return strings.TrimPrefix(forced, "s3::"), ok, err
}

func (g *S3Resolver) Resolve(name string) (afero.Fs, string, error) {
	u, err := urlhelper.Parse(name)
	if err != nil {
		return nil, "", errors.Wrap(err, "parse url")
	}
	region, bucket, path, _, creds, err := g.parseURL(u)
	if err != nil {
		return nil, "", errors.Wrap(err, "parse s3 url")
	}
	config := g.getAWSConfig(region, u, creds)
	if endpoint := os.Getenv("AWS_ENDPOINT"); endpoint != "" {
		config = config.WithEndpoint(endpoint)
	}
	sess := session.New(config)
	client := s3.New(sess)
	return aferoS3.NewFs(bucket, client), path, nil
}

func (g *S3Resolver) parseURL(u *url.URL) (region, bucket, path, version string, creds *credentials.Credentials, err error) {
	// This just check whether we are dealing with S3 or
	// any other S3 compliant service. S3 has a predictable
	// url as others do not
	if strings.Contains(u.Host, "amazonaws.com") {
		// Expected host style: s3.amazonaws.com. They always have 3 parts,
		// although the first may differ if we're accessing a specific region.
		hostParts := strings.Split(u.Host, ".")
		if len(hostParts) != 3 {
			err = fmt.Errorf("URL is not a valid S3 URL")
			return
		}

		// Parse the region out of the first part of the host
		region = strings.TrimPrefix(strings.TrimPrefix(hostParts[0], "s3-"), "s3")
		if region == "" {
			region = "us-east-1"
		}

		pathParts := strings.SplitN(u.Path, "/", 3)
		if len(pathParts) != 3 {
			err = fmt.Errorf("URL is not a valid S3 URL")
			return
		}

		bucket = pathParts[1]
		path = pathParts[2]
		version = u.Query().Get("version")

	} else {
		pathParts := strings.SplitN(u.Path, "/", 3)
		if len(pathParts) != 3 {
			err = fmt.Errorf("URL is not a valid S3 complaint URL")
			return
		}
		bucket = pathParts[1]
		path = pathParts[2]
		version = u.Query().Get("version")
		region = u.Query().Get("region")
		if region == "" {
			region = "us-east-1"
		}
	}

	_, hasAwsID := u.Query()["aws_access_key_id"]
	_, hasAwsSecret := u.Query()["aws_access_key_secret"]
	_, hasAwsToken := u.Query()["aws_access_token"]
	if hasAwsID || hasAwsSecret || hasAwsToken {
		creds = credentials.NewStaticCredentials(
			u.Query().Get("aws_access_key_id"),
			u.Query().Get("aws_access_key_secret"),
			u.Query().Get("aws_access_token"),
		)
	}

	return
}

func (g *S3Resolver) getAWSConfig(region string, url *url.URL, creds *credentials.Credentials) *aws.Config {
	conf := &aws.Config{}
	if creds == nil {
		// Grab the metadata URL
		metadataURL := os.Getenv("AWS_METADATA_URL")
		if metadataURL == "" {
			metadataURL = "http://169.254.169.254:80/latest"
		}

		creds = credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
				&ec2rolecreds.EC2RoleProvider{
					Client: ec2metadata.New(session.New(&aws.Config{
						Endpoint: aws.String(metadataURL),
					})),
				},
			})
	}

	if creds != nil {
		conf.Endpoint = &url.Host
		conf.S3ForcePathStyle = aws.Bool(true)
		if url.Scheme == "http" {
			conf.DisableSSL = aws.Bool(true)
		}
	}

	conf.Credentials = creds
	if region != "" {
		conf.Region = aws.String(region)
	}

	return conf
}
