package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

var (
	storageClient         *armstorage.ClientFactory
	storageAccountsClient *armstorage.AccountsClient
	subscriptionId        = "2be1e3f8-aa3d-455e-9c7c-974ab3077163"
)

type resourceId struct {
	SubscriptionID    string
	ResourceGroupName string
	ProviderNamespace string
	ResourceType      string
	ResourceName      string
}

func main() {
	cred, client, err := auth(subscriptionId)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	_, err = client.Get(context.TODO(), subscriptionId, nil)
	if err != nil {
		// TODO: handle error
		println("Error: ", err)
	}

	subscriptions := make([]*armsubscription.Subscription, 0)
	fmt.Println("Subscription: ", &subscriptions)

	storageClient, err = armstorage.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	storageAccountsClient = storageClient.NewAccountsClient()

	armClient, err := armresources.NewClient(subscriptionId, cred, nil)
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
				resourceId := *rg.ID
				resources[resourceId]++
				fmt.Printf("\r\nStorageAccount. ID: %s, ResourceGroup: %s\r\n", *rg.ID, *rg.Name)

				// extract the resource name
				resource, err := parseResourceID(resourceId)
				if err != nil {
					fmt.Errorf("Error parsing resource Id: ", err)
				}

				fmt.Println("Resource name: ", resource.ResourceName)
				// Get the storage account properties
				properties, err := getStorageAccountProperties(resourceId, resource.ResourceName)
				if err != nil {
					fmt.Errorf("Error with SA Properties: ", err)
				}

				fmt.Println("test")
				fmt.Println(properties.Name)
				// Get tags applied to storage account
				// GetResourceTags(subscriptionId, resourceID, cred)
			}
		}
	}
}

func getStorageAccountProperties(resourceGroup string, storageAccountName string) (*armstorage.Account, error) {
	storageAccountProperties, err := storageAccountsClient.GetProperties(context.TODO(), resourceGroup, storageAccountName, "2024-08-01")
	if err != nil {
		fmt.Println("SA Error: ", err)
		return &storageAccountProperties.Account, err
	}

	return &storageAccountProperties.Account, nil
}

func getResourceTags(subscriptionId string, resourceId string, credential azcore.TokenCredential) map[string]string {
	kvMap := make(map[string]string)
	tags, err := armresources.NewTagsClient(subscriptionId, credential, nil)
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

func auth(subscriptionId string) (azcore.TokenCredential, *armsubscription.SubscriptionsClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
		println("Error: ", err)
	}

	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		// TODO: handle error
	}

	return cred, client, err
}

func parseResourceID(resourceID string) (*resourceId, error) {
	parts := strings.Split(resourceID, "/")
	if len(parts) < 9 {
		return nil, fmt.Errorf("invalid resource ID format")
	}

	// Extract relevant parts based on Azure resource ID structure
	parsedId := &resourceId{
		SubscriptionID:    parts[2],
		ResourceGroupName: parts[4],
		ProviderNamespace: parts[6],
		ResourceType:      parts[7],
		ResourceName:      parts[8],
	}

	return parsedId, nil
}
