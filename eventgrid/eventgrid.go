package eventgrid

import (
	"context"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/helpers"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/iam"
	mgmt "github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2018-01-01/eventgrid"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/to"
)

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
