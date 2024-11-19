package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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
			if *rg.Type == "Microsoft.Storage/storageAccounts" {
				resourceId := *rg.ID
				resources[resourceId]++
				fmt.Printf("\r\nStorageAccount. ID: %s, ResourceGroup: %s\r\n", *rg.ID, *rg.Name)

				// Get tags applied to storage account
				// getResourceTags(subscriptionId, resourceId, cred) // works
				// createResourceTags(subscriptionId, resourceId, "CreateNewKey", "CreateNewValue", cred) // working
				deleteResourceTags(subscriptionId, resourceId, "test", cred)
			}
		}
	}
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

func createResourceTags(subscriptionId string, resourceId string, newKey string, newValue string, credential azcore.TokenCredential) map[string]string {
	kvMap := make(map[string]string)
	tags, err := armresources.NewTagsClient(subscriptionId, credential, nil)
	if err != nil {
		fmt.Errorf("Error creating tags client: %s", err)
	}

	resourceTags, err := tags.GetAtScope(context.TODO(), resourceId, nil)
	if err != nil {
		fmt.Errorf("Error getting tags: %s", err)
	}

	existingTags := map[string]*string{}
	if resourceTags.Properties != nil {
		existingTags = resourceTags.Properties.Tags
	}

	// Add our new tags
	existingTags[newKey] = &newValue

	// Merge existing tags with this
	tagsRequest := armresources.TagsResource{
		Properties: &armresources.Tags{
			Tags: existingTags,
		},
	}

	update, err := tags.CreateOrUpdateAtScope(
		context.TODO(),
		resourceId,
		tagsRequest,
		nil,
	)
	if err != nil {
		fmt.Println("Error creating tags: ", err)
	}

	fmt.Println(update)

	if kvMap != nil {
		return kvMap
	} else {
		fmt.Errorf("Map is nil")
		return nil
	}
}

func updateResourceTags(subscriptionId string, resourceId string, existingTags map[string]string, newKey string, newValue string, credential azcore.TokenCredential) map[string]string {
	kvMap := make(map[string]string)
	// tags, err := armresources.NewTagsClient(subscriptionId, credential, nil)
	// if err != nil {
	// 	fmt.Errorf("Error creating tags client: %s", err)
	// }

	// resourceTags, err := tags.GetAtScope(context.TODO(), resourceId, nil)
	// if err != nil {
	// 	fmt.Errorf("Error getting tags: %s", err)
	// }

	// // for key, value := range resourceTags.Properties.Tags {
	// // 	fmt.Println("[ExistingTags] Key:", key, "Value:", *value)
	// // 	kvMap[key] = *value
	// // 	for key, value := range existingTags {
	// // 		if key == existingTags[]
	// // 		tags.CreateOrUpdateValue(context.TODO(), existingTags[newKey], newValue, nil)
	// // 	}
	// }

	// // for key, value := range resourceTags.Properties.Tags {
	// // 	fmt.Println("[UpdatedTags] Key:", key, "Value:", *value)
	// // }

	if kvMap != nil {
		return kvMap
	} else {
		fmt.Errorf("Map is nil")
		return nil
	}
}

func deleteResourceTags(subscriptionId string, resourceId string, key string, credential azcore.TokenCredential) map[string]string {
	kvMap := make(map[string]string)
	tags, err := armresources.NewTagsClient(subscriptionId, credential, nil)
	if err != nil {
		fmt.Errorf("Error creating tags client: %s", err)
	}

	resourceTags, err := tags.GetAtScope(context.TODO(), resourceId, nil)
	if err != nil {
		fmt.Errorf("Error getting tags: %s", err)
	}

	existingTags := map[string]*string{}
	if resourceTags.Properties != nil {
		existingTags = resourceTags.Properties.Tags
	}

	// delete tag from our map
	// TODO: Error checking
	delete(existingTags, key)
	// del, err := delete(existingTags, key)
	// if err != nil {
	// 	fmt.Println("Error removing ", key, " Error: ", err)
	// }
	// fmt.Println(del)

	// Merge existing tags with this
	tagsRequest := armresources.TagsResource{
		Properties: &armresources.Tags{
			Tags: existingTags,
		},
	}

	update, err := tags.CreateOrUpdateAtScope(
		context.TODO(),
		resourceId,
		tagsRequest,
		nil,
	)
	if err != nil {
		fmt.Println("Error creating tags: ", err)
	}

	fmt.Println(update)

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
