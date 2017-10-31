package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	CREDENTIALS_TEMPLATES = `
		nexus_host = "{{ .Host }}"
		nexus_username = "{{ .Username }}"
		nexus_password = "{{ .Password }}"
		nexus_repository = "{{ .Repository }}"
	`
)

func main() {
	r := NewRegistry()

	app := cli.NewApp()
	app.Name = "Registry CLI"
	app.Usage = "Manage Docker Registries"
	app.Version = "1.0.0-beta"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Mohamed Labouardy",
			Email: "mohamed@labouardy.com",
		},
	}
	app.Commands = []cli.Command{
		{
			Name: "configure",
			Action: func(c *cli.Context) error {
				var hostname, repository, username, password string
				fmt.Print("Enter Nexus Host: ")
				fmt.Scan(&hostname)
				fmt.Print("Enter Nexus Repository Name: ")
				fmt.Scan(&repository)
				fmt.Print("Enter Nexus Username: ")
				fmt.Scan(&username)
				fmt.Print("Enter Nexus Password: ")
				fmt.Scan(&password)

				data := struct {
					Host       string
					Username   string
					Password   string
					Repository string
				}{
					hostname,
					username,
					password,
					repository,
				}

				tmpl, err := template.New(".credentials").Parse(CREDENTIALS_TEMPLATES)
				if err != nil {
					log.Fatal(err)
				}

				f, err := os.Create("swagger.yaml")
				if err != nil {
					log.Fatal(err)
				}

				err = tmpl.Execute(f, data)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name: "image",
			Subcommands: []cli.Command{
				{
					Name: "ls",
					Action: func(c *cli.Context) error {
						images := r.ListImages()
						for _, image := range images {
							fmt.Println(image)
						}
						fmt.Printf("Total images: %d", len(images))
						return nil
					},
				},
				{
					Name: "tags",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "List tags by image name",
						},
					},
					Action: func(c *cli.Context) error {
						if c.String("name") == "" {
							return cli.NewExitError("image name is required", 1)
						}
						tags := r.ListTagsByImage(c.String("name"))
						for _, tag := range tags {
							fmt.Println(tag)
						}
						fmt.Printf("There are %d images for %s", len(tags), c.String("name"))
						return nil
					},
				},
				{
					Name: "info",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name, n",
						},
						cli.StringFlag{
							Name: "tag, t",
						},
					},
					Action: func(c *cli.Context) error {
						if c.String("name") == "" || c.String("tag") == "" {
							return cli.NewExitError("image name & tag are required", 1)
						}
						manifest := r.ImageManifest(c.String("name"), c.String("tag"))
						fmt.Printf("Image: %s:%s\n", c.String("name"), c.String("tag"))
						fmt.Printf("Size: %d\n", manifest.Config.Size)
						fmt.Println("Layers:")
						for _, layer := range manifest.Layers {
							fmt.Printf("\t%s\t%d\n", layer.Digest, layer.Size)
						}
						return nil
					},
				},
				{
					Name: "delete",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name, n",
						},
						cli.StringFlag{
							Name: "tag, t",
						},
					},
					Action: func(c *cli.Context) error {
						if c.String("name") == "" || c.String("tag") == "" {
							return cli.NewExitError("image name & tag are required", 1)
						}
						r.DeleteImageByTag(c.String("name"), c.String("tag"))
						return nil
					},
				},
			},
		},
	}
	app.Run(os.Args)

}

// $ registry configure
// $ registry image ls
// $ registry image tags -name
// $ registry image info -name -tag
// $ registry image delete -name -tag
// $ registry image delete -name -keep 4
// $ registry image delete -keep 4
