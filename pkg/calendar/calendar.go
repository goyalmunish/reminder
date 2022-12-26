package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"net/http"
	"os"
	"time"

	"github.com/goyalmunish/reminder/pkg/logger"
	"github.com/goyalmunish/reminder/pkg/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gc "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const TitlePrefix string = "Reminder | "

func ConstructEvent(title string, description string, start time.Time, timezoneIANA string) (*gc.Event, error) {
	title = fmt.Sprintf("%s%s", TitlePrefix, title)
	startRFC3339 := &gc.EventDateTime{
		DateTime: start.Format(time.RFC3339),
		TimeZone: timezoneIANA,
	}
	endRFC3339 := &gc.EventDateTime{
		DateTime: start.Add(time.Duration(1 * 60 * 60 * time.Second)).Format(time.RFC3339),
		TimeZone: timezoneIANA,
	}
	source := &gc.EventSource{
		Title: "reminder",
		Url:   "https://github.com/goyalmunish/reminder",
	}
	rem := &gc.EventReminders{
		Overrides:  []*gc.EventReminder{},
		UseDefault: true,
	}
	recurrence := []string{"RRULE:FREQ=YEARLY"}
	event := &gc.Event{
		// ICalUID
		// Id
		// Created
		// Updated
		Summary:      title,
		Description:  description,
		Start:        startRFC3339,
		End:          endRFC3339,
		Recurrence:   recurrence,
		ColorId:      "10", // "Basil" color
		Reminders:    rem,
		EventType:    "default",
		Source:       source,
		Status:       "confirmed",
		Transparency: "transparent",
		Visibility:   "default",
	}
	return event, nil
}

// EventDetails returns overall event details.
func EventDetails(ctx context.Context, events *gc.Events) string {
	localTime := func(events gc.Events) string {
		value, err := utils.StrToTime(events.Updated, events.TimeZone)
		if err != nil {
			logger.Fatal(ctx, err)
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

// Get Calendar Service.
func GetCalendarService(ctx context.Context, options *Options) (*gc.Service, error) {
	credFile := options.CredentialFile
	b, err := os.ReadFile(utils.TryConvertTildaBasedPath(credFile))
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to read client secret file %q: %v", credFile, err))
	}
	logger.Info(ctx, fmt.Sprintf("Read client secret file %q", err))

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gc.CalendarEventsScope)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to parse client secret file to config: %v", err))
	}
	client := getClient(ctx, config, options)

	srv, err := gc.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, config *oauth2.Config, options *Options) *http.Client {
	// The file calendar_token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	tokenFile := options.TokenFile
	tok, err := tokenFromFile(ctx, utils.TryConvertTildaBasedPath(tokenFile))
	if err != nil {
		// Token file doesn't exist; envoke the authentication process
		// to generate one.
		tok = getTokenFromWeb(ctx, config)
		saveToken(ctx, tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(ctx context.Context, file string) (*oauth2.Token, error) {
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
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code and hit ENTER: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to read authorization code: %v", err))
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to retrieve token from web: %v", err))
	}
	return tok
}

// Saves a token to a file path.
func saveToken(ctx context.Context, path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(utils.TryConvertTildaBasedPath(path), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to cache oauth token: %v", err))
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to encode token: %v", err))
	}
}
