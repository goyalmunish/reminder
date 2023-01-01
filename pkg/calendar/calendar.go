package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"net/http"
	"os"

	"github.com/goyalmunish/reminder/pkg/logger"
	"github.com/goyalmunish/reminder/pkg/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gc "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const TitlePrefix string = "[reminder] "

// EventsDetails returns overall event details.
func EventsDetails(events *gc.Events) string {
	localTime := func(events gc.Events) string {
		value, err := utils.StrToTime(events.Updated, events.TimeZone)
		if err != nil {
			logger.Fatal(err)
		}
		return value.String()
	}
	reportTemplate := `
Calendar details:
  - Summary:  {{.Summary}}
  - TimeZone: {{.TimeZone}}
  - Updated:  {{. | timeInLocation}}
`
	funcMap := template.FuncMap{
		"timeInLocation": localTime,
	}
	return utils.TemplateResult(reportTemplate, funcMap, *events)
}

func EventString(event *gc.Event) string {
	details := []string{}
	details = append(details, event.Summary)
	details = append(details, event.Start.DateTime)
	details = append(details, event.Recurrence...)
	return strings.Join(details, " | ")
}

// Get Calendar Service.
func GetCalendarService(options *Options) (*gc.Service, error) {
	credFile := options.CredentialFile
	b, err := os.ReadFile(utils.TryConvertTildaBasedPath(credFile))
	if err != nil {
		return nil, fmt.Errorf("Couldn't read the client secret file %q; Refer instructions on https://github.com/goyalmunish/reminder#setting-up-the-environment-for-google-calendar-sync; Underneath error: %v", credFile, err)
	}
	logger.Info(fmt.Sprintf("Read client secret file %q.", credFile))

	// If modifying these scopes, delete your previously saved token file.
	config, err := google.ConfigFromJSON(b, gc.CalendarEventsScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config; If you changed the scope, then deleted your current %q token file and try again; Underneath error: %v", options.TokenFile, err)
	}

	client, err := getClient(config, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	srv, err := gc.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, options *Options) (*http.Client, error) {
	// The file calendar_token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	tokenFile := options.TokenFile
	tok, err := tokenFromFile(utils.TryConvertTildaBasedPath(tokenFile))
	if err != nil {
		logger.Warn(fmt.Sprintf("Token file doesn't exist; envoking the authentication process to generate one at %q.", tokenFile))
		tok, err := getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		saveToken(tokenFile, tok)
		logger.Info(fmt.Sprintf("Saved the token file %q.", tokenFile))
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser, find the authorization code in the URL and type it here and hit ENTER: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(utils.TryConvertTildaBasedPath(path), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to cache oauth token: %v", err))
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Unable to encode token: %v", err))
	}
}
