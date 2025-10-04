package main

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resource"
)

func main() {
	ctx := context.Background()
	subscription := "<subscription-id>"
	provider := "Microsoft.PolicyInsights"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Failed to create credential: %v", err)
	}

	client := resource.NewProvidersClient(subscription, cred, nil)
	_, err = client.Register(ctx, provider, nil)
	if err != nil {
		log.Fatalf("Failed to register provider: %v", err)
	}
}
