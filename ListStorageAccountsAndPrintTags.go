package main

import (
	"context"
	"fmt"
	"log"

	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

const subscriptionID = "2be1e3f8-aa3d-455e-9c7c-974ab3077163"

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
		println("Error: ", err)
	}

	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		// TODO: handle error
	}
	_, err = client.Get(context.TODO(), subscriptionID, nil)
	if err != nil {
		// TODO: handle error
		println("Error: ", err)
	}

	subscriptions := make([]*armsubscription.Subscription, 0)

	print(subscriptions)

	armClient, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	resources := make(map[string]int)
	pager := armClient.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			fmt.Errorf("failed to get resources: %v", err)
		}

		for _, rg := range page.Value {
			// fmt.Printf("\r\n ResourceID: %s, \r\nResource group: %s, \r\nLocation %s\r\n", *rg.ID, *rg.Name, *rg.Location)
			if *rg.Type == "Microsoft.Storage/storageAccounts" {
				resourceID := *rg.ID
				resources[resourceID]++
				fmt.Printf("\r\nStorageAccount. ID: %s, ResourceGroup: %s\r\n", *rg.ID, *rg.Name)
				GetResourceTags(subscriptionID, resourceID, cred)
			}
		}
	}
}

func GetResourceTags(subscriptionId string, resourceId string, credential azcore.TokenCredential) map[string]string {
	kvMap := make(map[string]string)
	tags, err := armresources.NewTagsClient(subscriptionID, credential, nil)
	if err != nil {
		fmt.Errorf("Error creating tags client: %s", err)
	}
	resourceTags, err := tags.GetAtScope(context.TODO(), resourceId, nil)
	if err != nil {
		fmt.Errorf("Error getting tags: %s", err)
	}
	for key, value := range resourceTags.Properties.Tags {
		fmt.Println("[Tags] Key:", key, "Value:", *value)
		kvMap[key] = *value
	}

	if kvMap != nil {
		return kvMap
	} else {
		fmt.Errorf("Map is nil")
		return nil
	}
}
