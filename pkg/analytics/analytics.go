package analytics

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/observer"
)

const (
	analyticsEndpoint = "https://www.google-analytics.com/collect"

	// Google Analytics hit types

	// PageViewHitType measure page views
	PageViewHitType = "pageview"
	// ScreenViewHitType measure sceen views in apps
	ScreenViewHitType = "screenview"
	// EventHitType measure any kind of event
	EventHitType = "event"
	// TransactionHitType measure e-commerce transactions
	TransactionHitType = "transaction"
	// ItemHitType measure items related to e-commerce transactions
	ItemHitType = "item"
	// SocialHitType measure activities on social media
	SocialHitType = "social"
)

type (
	// ReportFunc is the func interface used implement various 'uploaders'
	ReportFunc func(*Event)

	// Event  contains, well, events
	Event struct {
		Type      string // hit type
		Timestamp int64  // timestamp the event was created
		Client    *ClientRequest
		Values    map[string]string // k/v pairs of event data
	}

	// ClientRequest identifies a user, client request etc
	ClientRequest struct {
		IP          string
		Host        string
		Request     string
		UserID      string
		Fingerprint string
		UserAgent   string
	}
)

var (
	measurementID string
	appID         string
	events        chan (*Event)
	errorClient   *observer.Client
	reporters     []ReportFunc
)

func init() {
	projectID := env.GetString("PROJECT_ID", "")
	if projectID == "" {
		log.Fatal("Missing variable 'PROJECT_ID'")
	}
	// measurementID is the new Google Analytics ID
	measurementID = env.GetString("MEASUREMENT_ID", "UA-xxxxxxxx-x")
	appID = env.GetString("SERVICE_NAME", "backend")

	cl, err := observer.NewClient(context.Background(), projectID, appID)
	if err != nil {
		log.Fatal(err)
	}
	errorClient = cl

	AddReporter(LogReporter)

	events = make(chan *Event, 20)
	go analyticsUploader()
}

// AddReporter appends a ReportFunc to the array of uploaders that are called in turn
func AddReporter(r ReportFunc) {
	reporters = append(reporters, r)
}

// LogReporter is the most basic reporter func
func LogReporter(e *Event) {
	log.Println(e)
}

// CreateClientRequest returns a ClientRequest struct
func CreateClientRequest(request *http.Request) *ClientRequest {
	ip := request.RemoteAddr
	userAgent := request.UserAgent()

	return &ClientRequest{
		IP:          ip,
		Host:        request.Host,
		Request:     request.URL.Path,
		Fingerprint: util.Fingerprint(userAgent + ip),
		UserAgent:   userAgent,
	}
}

// TrackEvent post an event to analytics
func TrackEvent(request *http.Request, client, category, action, label string, value int) {
	v := make(map[string]string)
	v["ec"] = category
	v["ea"] = action
	v["el"] = label
	v["ev"] = strconv.FormatInt(int64(value), 10)

	e := &Event{
		Type:      EventHitType,
		Timestamp: util.Timestamp(),
		Client:    CreateClientRequest(request),
		Values:    v,
	}
	e.Client.UserID = client

	// push to the channel
	events <- e
}

// TrackPageView posts a pageview event to analytics
func TrackPageView(request *http.Request, client string) {
	e := &Event{
		Type:      PageViewHitType,
		Timestamp: util.Timestamp(),
		Client:    CreateClientRequest(request),
	}
	e.Client.UserID = client

	// push to the channel
	events <- e
}

func analyticsUploader() {
	for {
		e := <-events

		// iterate over all the reporters and send the event there
		for _, report := range reporters {
			report(e)
		}
	}
}

/*
func trackEventOld(request *http.Request, category, action, label string, value int) {
	v := make(map[string]string)

	// see https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#events
	v["ec"] = category
	v["ea"] = action
	v["el"] = label
	v["ev"] = strconv.FormatInt(int64(value), 10)

	PostToGoogleAnalytics(request, &Event{Type: EventHitType, Timestamp: util.Timestamp(), Values: v})
}

// TrackPageView posts a pageview event to analytics
func TrackPageView(request *http.Request) {
	PostToGoogleAnalytics(request, &Event{Type: PageViewHitType, Timestamp: util.Timestamp()})
}

// PostToGoogleAnalytics send the values to Google Analytics
func PostToGoogleAnalytics(request *http.Request, e *Event) {

	ip := request.RemoteAddr
	userAgent := request.UserAgent()
	uid := util.Fingerprint(userAgent + ip)

	// empty k/v is allowed
	if e.Values == nil {
		e.Values = make(map[string]string)
	}

	// enrich the event with some basics
	// see https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#content

	// e.Values["v"] = "1" this will be added just before posting
	e.Values["tid"] = measurementID
	e.Values["ds"] = appID
	e.Values["uid"] = uid
	e.Values["uip"] = ip
	e.Values["ua"] = url.QueryEscape(userAgent)
	e.Values["dh"] = url.QueryEscape(request.Host)
	e.Values["dp"] = url.QueryEscape(request.URL.Path)
	e.Values["npa"] = "1" // Disabling Advertising Personalization

	events <- e
}

func googleAnalyticsUploader() {
	for {
		e := <-events

		values := url.Values{
			"v": {"1"},
		}
		for k, v := range e.Values {
			vv := make([]string, 1)
			vv[0] = url.QueryEscape(v)
			values[k] = vv
		}

		resp, err := http.PostForm(analyticsEndpoint, values)
		if err != nil {
			errorClient.ReportError(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			errorClient.ReportError(fmt.Errorf("Google Analytics returned '%d'", resp.StatusCode))
			return
		}
	}
}
*/
