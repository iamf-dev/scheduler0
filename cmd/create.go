package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"scheduler0/config"
	"scheduler0/headers"
	"scheduler0/models"
	"scheduler0/secrets"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a resource like credential, projects or jobs",
	Long: `
Use this 

Usage:
	create credential
`,
}

var credentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "Creates a new credential",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr, "[cmd] ", log.LstdFlags)

		configs := config.Configurations(logger)
		credentials := secrets.GetSecrets(logger)

		credentialModel := models.CredentialModel{}
		data, err := credentialModel.ToJSON()
		if err != nil {
			logger.Fatalln(err)
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				req.Method = http.MethodPost
				body := bytes.NewReader(data)
				rc := ioutil.NopCloser(body)
				req.Body = rc
				req.SetBasicAuth(credentials.AuthUsername, credentials.AuthPassword)
				req.Header.Add(headers.PeerHeader, "cmd")
				req.Header.Add("Content-Type", "application/json")

				if len(via) > 5 {
					logger.Fatalln("Too many redirects")
				}

				return nil
			},
		}
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s://%s:%s/credentials", configs.Protocol, configs.Host, configs.Port),
			bytes.NewReader(data),
		)

		req.SetBasicAuth(credentials.AuthUsername, credentials.AuthPassword)
		req.Header.Add(headers.PeerHeader, "cmd")
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			logger.Fatalln(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Fatalln(err)
		}

		if res.StatusCode != http.StatusCreated {
			logger.Fatalln("failed to create new credential:error:", string(body))
		} else {
			logger.Println("successfully created a credential ", string(body))
		}
	},
}

func init() {
	CreateCmd.AddCommand(credentialCmd)
}
