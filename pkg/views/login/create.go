package login

import (
	"errors"
	"net/url"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type LoginView struct {
	Server   string
	Username string
	Password string
	Name     string
}

func CreateView(loginView *LoginView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server").
				Value(&loginView.Server).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("server cannot be empty or only spaces")
					}
					formattedUrl := utils.FormatUrl(str)
					if _, err := url.ParseRequestURI(formattedUrl); err != nil {
						return errors.New("enter the correct server format")
					}
					return nil
				}),
			huh.NewInput().
				Title("User Name").
				Value(&loginView.Username).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("username cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&loginView.Password).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("password cannot be empty")
					}
					isVaild, err := utils.ValidatePassword(str)
					if !isVaild {
						return err
					}
					return nil
				}),
			huh.NewInput().
				Title("Name of Credential").
				Value(&loginView.Name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("credential name cannot be empty or only spaces")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}

}
