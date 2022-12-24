package model

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/goyalmunish/reminder/pkg/logger"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var HomeDir string = os.Getenv("HOME")

const EnableCalendar bool = false

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	// The file calendar_token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first time.
	tokenFile := fmt.Sprintf("%s/%s", HomeDir, "calendar_token.json")
	tok, err := tokenFromFile(ctx, tokenFile)
	if err != nil {
		tok = getTokenFromWeb(ctx, config)
		saveToken(ctx, tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
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

// Saves a token to a file path.
func saveToken(ctx context.Context, path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to cache oauth token: %v", err))
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to encode token: %v", err))
	}
}

// Get Calendar Service.
func getCalendarService(ctx context.Context) (*calendar.Service, error) {
	credFile := fmt.Sprintf("%s/%s", HomeDir, "calendar_credentials.json")
	b, err := os.ReadFile(credFile)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to read client secret file: %v", err))
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to parse client secret file to config: %v", err))
	}
	client := getClient(ctx, config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

// Fetch upcoming calendar events.
func FetchCalendar(ctx context.Context) {
	if !EnableCalendar {
		logger.Warn(ctx, "Google Calendar is disabled.")
		return
	}
	srv, err := getCalendarService(ctx)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to retrieve Calendar client: %v", err))
	}
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Unable to retrieve next ten of the user's events: %v", err))
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}
