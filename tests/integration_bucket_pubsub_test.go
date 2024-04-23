package test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func create_empty_file(fileName string) {
	emptyContent := []byte("")
	err := os.WriteFile(fileName, emptyContent, 0644)
	if err != nil {
		fmt.Printf("Error creating empty file: %v\n", err)
		return
	}
}

func upload_file_to_bucket(ctx context.Context, client *storage.Client,
	fileName string, bucketName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Upload the file to the bucket
	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		fmt.Printf("Error uploading file to bucket: %v\n", err)
		return
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("Error closing bucket writer: %v\n", err)
		return
	}

}

func TestIntegrationBucketPubSub(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.

	projectID := "temporal-field-419719"
	subscriptionID := "test_integration_subscription"
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "../modules/bucket",
		Vars: map[string]interface{}{
			"project":           projectID, // Set the instance type variable
			"name":              "test_integration_bucket_123",
			"location":          "US",
			"topic_name":        "test_integration_bucket_topic",
			"subscription_name": "test_integration_subscription",
		},
	})

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error creating storage client: %v\n", err)
		return
	}
	defer client.Close()

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// bucket name
	output_bucket_name := terraform.Output(t, terraformOptions, "name")

	// create file in bucket
	fileName := "empty_file.txt"
	create_empty_file(fileName)
	upload_file_to_bucket(ctx, client, fileName, output_bucket_name)

	// Create a new Pub/Sub client using default authentication
	pubsub_client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}

	// Get a subscription handle
	sub := pubsub_client.Subscription(subscriptionID)

	// Create a context with a timeout to stop the message pulling after a certain duration
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Receive messages
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// Handle the received message
		fmt.Printf("Got message: %s\n", string(msg.Data))

		// Acknowledge the message to indicate it has been processed
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("Error receiving message: %v", err)
	}

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	pattern := regexp.MustCompile(`-\p{L}{4}$`)

	assert.True(t, pattern.MatchString(output_bucket_name))
}

/*
gcloud projects add-iam-policy-binding temporal-field-419719 \
    --member=serviceAccount:service-287505208604@gs-project-accounts.iam.gserviceaccount.com \
    --role=roles/pubsub.publisher
*/
