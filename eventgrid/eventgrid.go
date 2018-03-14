package eventgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/helpers"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/iam"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	mgmt "github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2018-01-01/eventgrid"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/to"
)

func getSubscriptionClient() (client mgmt.EventSubscriptionsClient, err error) {
	var token adal.OAuthTokenProvider
	token, err = iam.GetResourceManagementToken(iam.AuthGrantType())
	if err != nil {
		return
	}

	client = mgmt.NewEventSubscriptionsClient(helpers.SubscriptionID())
	client.Authorizer = autorest.NewBearerAuthorizer(token)
	client.AddToUserAgent(helpers.UserAgent())
	return
}

func CreateTopic(ctx context.Context, name string) (created mgmt.Topic, err error) {
	client, err := getTopicsClient()
	if err != nil {
		return
	}

	var fut mgmt.TopicsCreateOrUpdateFuture
	fut, err = client.CreateOrUpdate(ctx, helpers.ResourceGroupName(), name, mgmt.Topic{
		Name:     to.StringPtr(name),
		Location: to.StringPtr(helpers.Location()),
	})
	if err != nil {
		return
	}

	err = fut.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return
	}

	created, err = fut.Result(client)
	return
}

func CreateEndpoint(ctx context.Context, address string) {

}

// MockTopic creates a local message dispatcher, loosely emulating an Azure hosted
// EventGrid for the sake of the examples in this package. If you ever find yourself
// wanting to take a dependency on this, you should really ask yourself if you:
//    A) Need a real instance of EventGrid.
//    B) Can do something lighter-weight with channels directly.
type MockTopic struct {
	sync.RWMutex
	subscribers map[mockTopicSubscribers]struct{}
}

type mockTopicSubscribers struct {
	callback func(context.Context, eventgrid.Event) int
	filter   *mgmt.EventSubscriptionFilter
}

func (mt *MockTopic) ListenAndServe(addr string) error {
	http.HandleFunc("/", mt.handler)
	return http.ListenAndServe(addr, nil)
}

func (mt *MockTopic) handler(w http.ResponseWriter, req *http.Request) {
	const maxEventSize = 64 * 1024
	const maxEvents = 100
	const maxPayloadSize = maxEventSize * maxEvents

	lr := io.LimitReader(req.Body, maxPayloadSize)
	body, err := ioutil.ReadAll(lr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to read request")
		return
	}

	var payload []eventgrid.Event
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "request was not a JSON array of EventGrid events")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, event := range payload {
		for sub := range mt.subscribers {
			if applyFilter(sub.filter, event) {

			}
		}
	}
}

func applyFilter(filter *mgmt.EventSubscriptionFilter, event eventgrid.Event) bool {
	if filter.IncludedEventTypes != nil && event.EventType != nil {
		foundEventType := false
		for _, included := range *filter.IncludedEventTypes {
			if strings.EqualFold(included, *event.EventType) {
				foundEventType = true
				break
			}
		}
		if !foundEventType {
			return false
		}
	}

	var eventSubject 
}

func (mt *MockTopic) Endpoint() string {
}

// Subscribe mimicks the functionality of registering an Event Handler with a Topic via ARM operations.
func (mt *MockTopic) Subscribe(callback func(context.Context, eventgrid.Event) int, filter *mgmt.EventSubscriptionFilter) {
	mt.Lock()
	defer mt.Unlock()
}

func (mt *MockTopic) Publish(event eventgrid.Event) {
	mt.RLock()
	defer mt.RUnlock()

}

func getTopicsClient() (client mgmt.TopicsClient, err error) {
	var token adal.OAuthTokenProvider
	token, err = iam.GetResourceManagementToken(iam.AuthGrantType())
	if err != nil {
		return
	}

	client = mgmt.NewTopicsClient(helpers.SubscriptionID())
	client.Authorizer = autorest.NewBearerAuthorizer(token)
	client.AddToUserAgent(helpers.UserAgent())
	return
}
