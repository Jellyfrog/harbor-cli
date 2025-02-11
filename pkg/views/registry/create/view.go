package create

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// struct to hold registry options
type RegistryOption struct {
	ID   string
	Name string
}

func CreateRegistryView(createView *api.CreateRegView) {
	registries, _ := api.GetRegistryProviders()

	// Initialize a slice to hold registry options
	var registryOptions []RegistryOption

	// Iterate over registries to populate registryOptions
	for i, registry := range registries {
		registryOptions = append(registryOptions, RegistryOption{
			ID:   strconv.FormatInt(int64(i), 10),
			Name: registry,
		})
	}

	// Initialize a slice to hold select options
	var registrySelectOptions []huh.Option[string]

	// Iterate over registryOptions to populate registrySelectOptions
	for _, option := range registryOptions {
		registrySelectOptions = append(
			registrySelectOptions,
			huh.NewOption(option.Name, option.Name),
		)
	}

	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a Registry Provider").
				Value(&createView.Type).
				Options(registrySelectOptions...).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("registry provider cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("name cannot be empty or only spaces")
					}
					if isVaild := utils.ValidateRegistryName(str); !isVaild {
						return errors.New("please enter the correct name format")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
			huh.NewInput().
				Title("URL").
				Value(&createView.URL).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("url cannot be empty or only spaces")
					}
					formattedUrl := utils.FormatUrl(str)
					if _, err := url.ParseRequestURI(formattedUrl); err != nil {
						return errors.New("please enter the correct url format")
					}
					return nil
				}),
			huh.NewInput().
				Title("Access Key").
				Value(&createView.Credential.AccessKey),
			huh.NewInput().
				Title("Access Secret").
				Value(&createView.Credential.AccessSecret),
			huh.NewConfirm().
				Title("Verify Cert").
				Value(&createView.Insecure).
				Affirmative("yes").
				Negative("no"),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}
