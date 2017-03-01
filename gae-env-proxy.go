package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "gae-env-proxy"
	app.Usage = "Proxy environment variables to app.yml file"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "prefix, p",
			Usage: "prefix to use for filtering env variables",
		},
		cli.StringFlag{
			Name:  "input, i",
			Usage: "input app.yml file",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "output app.yml file",
		},
	}

	app.Action = func(c *cli.Context) error {
		inputFile := os.Stdin
		outputFile := os.Stdout
		envPrefix := c.String("prefix")

		if f, err := os.Open(c.String("input")); err == nil {
			inputFile = f
			defer inputFile.Close()
		}
		if f, err := os.OpenFile(c.String("output"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666); err == nil {
			inputFile = f
			defer outputFile.Close()
		}
		bytes, err := ioutil.ReadAll(inputFile)
		if err != nil {
			return err
		}
		appMap := make(map[interface{}]interface{})
		err = yaml.Unmarshal(bytes, &appMap)

		if _, ok := appMap["env_variables"]; !ok {
			appMap["env_variables"] = make(map[interface{}]interface{})
		}
		appEnv := appMap["env_variables"].(map[interface{}]interface{})
		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			if strings.HasPrefix(pair[0], envPrefix) {
				appEnv[pair[0]] = pair[1]
			}
		}

		bytes, err = yaml.Marshal(appMap)
		_, err = outputFile.Write(bytes)
		if err != nil {
			return err
		}
		return nil
	}

	app.Run(os.Args)
}
