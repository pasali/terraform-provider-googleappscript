package googleappscript

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/script/v1"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Config struct {
	script *script.Service
	drive  *drive.Service
}

func (c *Config) loadAndValidate(tokenFile string) error {
	log.Printf("[INFO] authenticating with local client")
	client, _ := getClient(tokenFile)
	client.Transport = logging.NewTransport("Google", client.Transport)
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s",
		runtime.GOOS, runtime.GOARCH, terraform.VersionString())

	scriptSvc, err := script.New(client)
	if err != nil {
		return nil
	}
	scriptSvc.UserAgent = userAgent
	c.script = scriptSvc

	driveSvc, err := drive.New(client)
	if err != nil {
		return nil
	}
	driveSvc.UserAgent = userAgent
	c.drive = driveSvc
	return nil
}

func getClient(tokenFile string) (*http.Client, error) {
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	tokenSource := oauth2.StaticTokenSource(token)
	return oauth2.NewClient(context.Background(), tokenSource), err
}
