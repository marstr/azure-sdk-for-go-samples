package eventgrid

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/marstr/guid"
	"net/http"
)

// MyEventSpecificData demonstrates how one can declare an aribtrary type to be
// sent and received as part of an EventGrid call.
type MyEventSpecificData struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 string `json:"field3"`
}

type extendedEventGridClient struct {
	eventgrid.BaseClient
	Token string `json:"aeg-sas-token"`
	Key   string `json:"aeg-sas-key"`
}

func ExampleListen() {
	http.ListenAndServe()
}

func ExampleBaseClient_PublishEvents() {
	myPayload := &MyEventSpecificData{
		Field1: "value1",
		Field2: "value2",
		Field3: "value3",
	}

	myEvent := eventgrid.Event{
		EventTime:   &date.Time{Time: time.Now()},
		EventType:   to.StringPtr("MyEventSpecificData"),
		Data:        myPayload,
		DataVersion: to.StringPtr("v1.0.0"),
		ID:          to.StringPtr(guid.NewGUID().String()),
		Subject:     to.StringPtr("S.P.Line"),
	}

	client := extendedEventGridClient{
		Key: "mbQrL514nXMgXDwMz1zPfkODqA17jAmcNgAuVrMu/Gc=",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	resp, err := client.PublishEvents(ctx, "caseyjones.westus2-1.eventgrid.azure.net", []eventgrid.Event{myEvent})
	if err != nil {
		return
	}
	fmt.Println(resp.Status)

	// Output: 200 OK
}

// PublishEvents shadows the `BaseClient` method of the same name. It does so in order to add the required authentication Header,
// one of `aeg-sas-token` or `aeg-sas-key`.
//
// Note: This function is a work-around for issue https://github.com/Azure/azure-rest-api-specs/issues/2849 is fixed.
func (client extendedEventGridClient) PublishEvents(ctx context.Context, topicHostname string, events []eventgrid.Event) (result autorest.Response, err error) {
	err = validation.Validate([]validation.Validation{
		{
			TargetValue: events,
			Constraints: []validation.Constraint{
				{Target: "events", Name: validation.Null, Rule: true, Chain: nil},
			},
		},
	})
	if err != nil {
		return result, validation.NewError("eventgrid.extendedEventGridClient", "PublishEvents", err.Error())
	}

	req, err := client.PublishEventsPreparer(ctx, topicHostname, events)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.extendedEventGridClient", "PublishEvents", nil, "Failure preparing request")
	}

	if client.Key != "" {
		req.Header.Add("aeg-sas-key", client.Key)
	}

	if client.Token != "" {
		req.Header.Add("aeg-sas-token", client.Token)
	}

	resp, err := client.PublishEventsSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "eventgrid.BaseClient", "PublishEvents", resp, "Failure sending request")
		return
	}

	result, err = client.PublishEventsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventgrid.BaseClient", "PublishEvents", resp, "Failure responding to request")
	}

	return
}
