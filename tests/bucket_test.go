package test

import (
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestBucketCreatedWith4CharacterSuffix(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "../modules/bucket",
		Vars: map[string]interface{}{
			"project":    "temporal-field-419719", // Set the instance type variable
			"name":       "test_bucket_123",
			"location":   "US",
			"topic_name": "test_bucket_topic",
		},
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	pattern := regexp.MustCompile(`-\p{L}{4}$`)

	output_name := terraform.Output(t, terraformOptions, "name")
	assert.True(t, pattern.MatchString(output_name))
}

/*
gcloud projects add-iam-policy-binding temporal-field-419719 \
    --member=serviceAccount:service-287505208604@gs-project-accounts.iam.gserviceaccount.com \
    --role=roles/pubsub.publisher
*/
