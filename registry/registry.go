package registry

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env"
)

const acceptHeader = "application/vnd.docker.distribution.manifest.v2+json"
const credentialsFile = ".credentials"

// Registry struct to access registry information
type Registry struct {
	Host       string `toml:"nexus_host" env:"NEXUS_CLI_HOST"`
	Username   string `toml:"nexus_username" env:"NEXUS_CLI_USERNAME"`
	Password   string `toml:"nexus_password" env:"NEXUS_CLI_PASSWORD"`
	Repository string `toml:"nexus_repository" env:"NEXUS_CLI_REPOSITORY"`
}

// Repositories struct containing a slice of images
type Repositories struct {
	Images []string `json:"repositories"`
}

// ImageTags struct containing a slice of all tags for a given docker image name
type ImageTags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// ImageManifest struct for docker image information on schema and layers
type ImageManifest struct {
	SchemaVersion int64       `json:"schemaVersion"`
	MediaType     string      `json:"mediaType"`
	Config        LayerInfo   `json:"config"`
	Layers        []LayerInfo `json:"layers"`
}

// LayerInfo struct for docker image meta information
type LayerInfo struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

// NewRegistry uses local .credentials file or environment variables to return a Registry struct
func NewRegistry() (Registry, error) {
	r := Registry{}
	toml.DecodeFile(credentialsFile, &r)

	// Parse environment variables by struct `env`-tags
	env.Parse(&r)

	if len(r.Host) == 0 {
		return r, fmt.Errorf("Problem reading host from configuration")
	} else if len(r.Username) == 0 {
		return r, fmt.Errorf("Problem reading username from configuration")
	} else if len(r.Password) == 0 {
		return r, fmt.Errorf("Problem reading password from configuration")
	} else if len(r.Repository) == 0 {
		return r, fmt.Errorf("Problem reading repository from configuration")
	}
	return r, nil
}

// ListImages returns image names as a slice of strings
func (r Registry) ListImages() ([]string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/_catalog", r.Host, r.Repository)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	var repositories Repositories
	json.NewDecoder(resp.Body).Decode(&repositories)

	return repositories.Images, nil
}

// ListTagsByImage expects an image name as string to return a slice of tage names
func (r Registry) ListTagsByImage(image string) ([]string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/tags/list", r.Host, r.Repository, image)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	var imageTags ImageTags
	json.NewDecoder(resp.Body).Decode(&imageTags)

	return imageTags.Tags, nil
}

// ImageManifest expects image name and tag as string to return an ImageManifest struct
func (r Registry) ImageManifest(image string, tag string) (ImageManifest, error) {
	var imageManifest ImageManifest
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return imageManifest, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return imageManifest, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return imageManifest, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	json.NewDecoder(resp.Body).Decode(&imageManifest)

	return imageManifest, nil

}

// DeleteImageByTag expects an image name and a tag to delete an image tag
func (r Registry) DeleteImageByTag(image string, tag string) error {
	sha, err := r.getImageSHA(image, tag)
	if err != nil {
		return err
	}
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, sha)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	fmt.Printf("%s:%s has been successful deleted\n", image, tag)

	return nil
}

func (r Registry) getImageSHA(image string, tag string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	return resp.Header.Get("docker-content-digest"), nil
}
