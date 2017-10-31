package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/urfave/cli"
)

const (
	CREDENTIALS_TEMPLATES = `# Nexus Credentials
nexus_host = "{{ .Host }}"
nexus_username = "{{ .Username }}"
nexus_password = "{{ .Password }}"
nexus_repository = "{{ .Repository }}"`
)

func main() {
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
				return setNexusCredentials()
			},
		},
		{
			Name: "image",
			Subcommands: []cli.Command{
				{
					Name: "ls",
					Action: func(c *cli.Context) error {
						return listImages()
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
						return listTagsByImage(c.String("name"))
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
						return showImageInfo(c.String("name"), c.String("tag"))
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
						return deleteImageByTag(c.String("name"), c.String("tag"))
					},
				},
			},
		},
	}
	app.Run(os.Args)
}

func setNexusCredentials() error {
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
		return cli.NewExitError(err.Error(), 1)
	}

	f, err := os.Create(".credentials")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func listImages() error {
	r, err := NewRegistry()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	images, err := r.ListImages()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	for _, image := range images {
		fmt.Println(image)
	}
	fmt.Printf("Total images: %d", len(images))
	return nil
}

func listTagsByImage(imgName string) error {
	r, err := NewRegistry()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if imgName == "" {
		return cli.NewExitError("image name is required", 1)
	}
	tags, err := r.ListTagsByImage(imgName)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	for _, tag := range tags {
		fmt.Println(tag)
	}
	fmt.Printf("There are %d images for %s", len(tags), imgName)
	return nil
}

func showImageInfo(imgName string, tag string) error {
	r, err := NewRegistry()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if imgName == "" || tag == "" {
		return cli.NewExitError("image name & tag are required", 1)
	}
	manifest, err := r.ImageManifest(imgName, tag)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Printf("Image: %s:%s\n", imgName, tag)
	fmt.Printf("Size: %d\n", manifest.Config.Size)
	fmt.Println("Layers:")
	for _, layer := range manifest.Layers {
		fmt.Printf("\t%s\t%d\n", layer.Digest, layer.Size)
	}
	return nil
}

func deleteImageByTag(imgName string, tag string) error {
	r, err := NewRegistry()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if imgName == "" || tag == "" {
		return cli.NewExitError("image name & tag are required", 1)
	}
	err = r.DeleteImageByTag(imgName, tag)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// $ registry configure
// $ registry image ls
// $ registry image tags -name
// $ registry image info -name -tag
// $ registry image delete -name -tag
// $ registry image delete -name -keep 4
// $ registry image delete -keep 4
