package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/moepi/nexus-cli/registry"
	"github.com/urfave/cli"
)

const (
	credentialsTemplates = `# Nexus Credentials
nexus_host = "{{ .Host }}"
nexus_username = "{{ .Username }}"
nexus_password = "{{ .Password }}"
nexus_repository = "{{ .Repository }}"`
)

func main() {
	app := cli.NewApp()
	app.Name = "Nexus CLI"
	app.Usage = "Manage Docker Private Registry on Nexus"
	app.Version = "1.0.0-beta-2"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Mohamed Labouardy",
			Email: "mohamed@labouardy.com",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "configure",
			Usage: "Configure Nexus Credentials",
			Action: func(c *cli.Context) error {
				return setNexusCredentials(c)
			},
		},
		{
			Name:  "image",
			Usage: "Manage Docker Images",
			Subcommands: []cli.Command{
				{
					Name:  "ls",
					Usage: "List all images in repository",
					Action: func(c *cli.Context) error {
						return listImages(c)
					},
				},
				{
					Name:  "tags",
					Usage: "Display all image tags",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "List tags by image name",
						},
						cli.StringSliceFlag{
							Name:  "expression, e",
							Usage: "Filter tags by regular expression",
						},
						cli.BoolFlag{
							Name:  "invert, v",
							Usage: "Invert filter results",
						},
					},
					Action: func(c *cli.Context) error {
						return listTagsByImage(c)
					},
				},
				{
					Name:  "info",
					Usage: "Show image details",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name, n",
						},
						cli.StringFlag{
							Name: "tag, t",
						},
					},
					Action: func(c *cli.Context) error {
						return showImageInfo(c)
					},
				},
				{
					Name:  "delete",
					Usage: "Delete images",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name, n",
						},
						cli.StringFlag{
							Name: "tag, t",
						},
						cli.StringFlag{
							Name: "keep, k",
						},
						cli.StringSliceFlag{
							Name:  "expression, e",
							Usage: "Filter tags by regular expression",
						},
						cli.BoolFlag{
							Name:  "invert, v",
							Usage: "Invert results filter expressions",
						},
					},
					Action: func(c *cli.Context) error {
						return deleteImages(c)
					},
				},
			},
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Wrong command %q !", command)
	}
	app.Run(os.Args)
}

func setNexusCredentials(c *cli.Context) error {
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

	tmpl, err := template.New(".credentials").Parse(credentialsTemplates)
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

func listImages(c *cli.Context) error {
	r, err := registry.NewRegistry()
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
	fmt.Printf("Total images: %d\n", len(images))
	return nil
}

func filterTagsByRegex(tags []string, expressions []string, invert bool) ([]string, error) {
	var retTags []string
	if len(expressions) == 0 {
		return tags, nil
	}
	for _, tag := range tags {
		tagMiss := false
		for _, expression := range expressions {
			var expressionBool = !invert
			if strings.HasPrefix(expression, "!") {
				expressionBool = invert
				expression = strings.Trim(expression, "!")
			}
			retVal, err := regexp.MatchString(expression, tag)
			if err != nil {
				return retTags, err
			}
			if retVal != expressionBool {
				tagMiss = true
				break
			}
		}
		// tag must match all expression, so continue with next tag on match
		if !tagMiss {
			retTags = append(retTags, tag)
		}
	}
	return retTags, nil
}

func listTagsByImage(c *cli.Context) error {
	var imgName = c.String("name")
	r, err := registry.NewRegistry()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if imgName == "" {
		cli.ShowSubcommandHelp(c)
	}
	tags, err := r.ListTagsByImage(imgName)

	// filter tags by expressions
	tags, err = filterTagsByRegex(tags, c.StringSlice("expression"), c.Bool("invert"))
	if err != nil {
		log.Fatal(err)
	}

	compareStringNumber := func(str1, str2 string) bool {
		return extractNumberFromString(str1) < extractNumberFromString(str2)
	}
	Compare(compareStringNumber).Sort(tags)

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	for _, tag := range tags {
		fmt.Println(tag)
	}
	fmt.Printf("There are %d images for %s\n", len(tags), imgName)
	return nil
}

func showImageInfo(c *cli.Context) error {
	var imgName = c.String("name")
	var tag = c.String("tag")
	r, err := registry.NewRegistry()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if imgName == "" || tag == "" {
		cli.ShowSubcommandHelp(c)
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

func deleteImages(c *cli.Context) error {
	var imgName = c.String("name")
	var tag = c.String("tag")
	var keep = c.Int("keep")
	var invert = c.Bool("invert")

	// Show help if no image name is present
	if imgName == "" {
		fmt.Fprintf(c.App.Writer, "You should specify the image name\n")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	r, err := registry.NewRegistry()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// if a specific tag is provided, ignore all other options
	if tag != "" {
		err = r.DeleteImageByTag(imgName, tag)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}

	// Get list of tags and filter them by all expressions provided
	tags, err := r.ListTagsByImage(imgName)
	tags, err = filterTagsByRegex(tags, c.StringSlice("expression"), invert)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Could not filter tags by regular expressions: %s\n", err)
		return err
	}

	// if no keep is specified, all flags are unset. Show help and exit.
	if c.IsSet("keep") == false && len(c.StringSlice("expression")) == 0 {
		fmt.Fprintf(c.App.Writer, "You should either specify use tag / filter expressions, or specify how many images you want to keep\n")
		cli.ShowSubcommandHelp(c)
		return fmt.Errorf("You should either specify use tag / filter expressions, or specify how many images you want to keep")
	}

	if len(tags) == 0 && !c.IsSet("keep") {
		fmt.Fprintf(c.App.Writer, "No images selected for deletion\n")
		return fmt.Errorf("No images selected for deletion")
	}

	// Remove images by using keep flag
	compareStringNumber := func(str1, str2 string) bool {
		return extractNumberFromString(str1) < extractNumberFromString(str2)
	}
	Compare(compareStringNumber).Sort(tags)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if len(tags) >= keep {
		for _, tag := range tags[:len(tags)-keep] {
			fmt.Printf("%s:%s image will be deleted ...\n", imgName, tag)
			err = r.DeleteImageByTag(imgName, tag)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}
	} else {
		fmt.Printf("Only %d images are available\n", len(tags))
	}
	return nil
}
