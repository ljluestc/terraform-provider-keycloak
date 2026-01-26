package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestCustomIdentityProviderMapperKeyTransformation(t *testing.T) {
	resource := resourceKeycloakCustomIdentityProviderMapper()
	resourceSchema := resource.Schema

	// Test Case 1: attribute_value -> attribute.value transformation
	t.Run("attribute_value transformation", func(t *testing.T) {
		rawData := map[string]interface{}{
			"realm":                    "my-realm",
			"name":                     "my-mapper",
			"identity_provider_alias":  "my-idp",
			"identity_provider_mapper": "hardcoded-attribute-idp-mapper",
			"extra_config": map[string]interface{}{
				"attribute":       "something",
				"attribute_value": "some-value",
			},
		}

		data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

		mapper := getCustomIdentityProviderMapperFromData(data)

		// Verify transformation happened
		if val, ok := mapper.Config.ExtraConfig["attribute.value"]; !ok || val != "some-value" {
			t.Errorf("Expected extra_config['attribute.value'] to be 'some-value', got %v", val)
		}
		if _, ok := mapper.Config.ExtraConfig["attribute_value"]; ok {
			t.Errorf("Expected extra_config['attribute_value'] to be removed")
		}

		// Verify reverse transformation
		// Simulate Keycloak returning attribute.value
		returnedMapper := &keycloak.CustomIdentityProviderMapper{
			Realm:                  "my-realm",
			Name:                   "my-mapper",
			IdentityProviderAlias:  "my-idp",
			IdentityProviderMapper: "hardcoded-attribute-idp-mapper",
			Config: &keycloak.CustomIdentityProviderMapperConfig{
				ExtraConfig: map[string]interface{}{
					"attribute":       "something",
					"attribute.value": "some-value",
				},
			},
		}

		// We need a fresh ResourceData with the original config to test setCustomIdentityProviderMapperData correctly
		dataForSet := schema.TestResourceDataRaw(t, resourceSchema, rawData)

		setCustomIdentityProviderMapperData(dataForSet, returnedMapper)

		extraConfig := dataForSet.Get("extra_config").(map[string]interface{})
		if val, ok := extraConfig["attribute_value"]; !ok || val != "some-value" {
			t.Errorf("Expected extra_config['attribute_value'] to be 'some-value', got %v", val)
		}
		if _, ok := extraConfig["attribute.value"]; ok {
			t.Errorf("Expected extra_config['attribute.value'] to be removed")
		}
	})

	// Test Case 2: claim_value -> claim.value transformation
	t.Run("claim_value transformation", func(t *testing.T) {
		rawData := map[string]interface{}{
			"realm":                    "my-realm",
			"name":                     "my-mapper",
			"identity_provider_alias":  "my-idp",
			"identity_provider_mapper": "hardcoded-claim-idp-mapper",
			"extra_config": map[string]interface{}{
				"claim":       "something",
				"claim_value": "some-value",
			},
		}

		data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

		mapper := getCustomIdentityProviderMapperFromData(data)

		// Verify transformation happened
		if val, ok := mapper.Config.ExtraConfig["claim.value"]; !ok || val != "some-value" {
			t.Errorf("Expected extra_config['claim.value'] to be 'some-value', got %v", val)
		}
		if _, ok := mapper.Config.ExtraConfig["claim_value"]; ok {
			t.Errorf("Expected extra_config['claim_value'] to be removed")
		}

		// Verify reverse transformation
		returnedMapper := &keycloak.CustomIdentityProviderMapper{
			Realm:                  "my-realm",
			Name:                   "my-mapper",
			IdentityProviderAlias:  "my-idp",
			IdentityProviderMapper: "hardcoded-claim-idp-mapper",
			Config: &keycloak.CustomIdentityProviderMapperConfig{
				ExtraConfig: map[string]interface{}{
					"claim":       "something",
					"claim.value": "some-value",
				},
			},
		}

		dataForSet := schema.TestResourceDataRaw(t, resourceSchema, rawData)
		setCustomIdentityProviderMapperData(dataForSet, returnedMapper)

		extraConfig := dataForSet.Get("extra_config").(map[string]interface{})
		if val, ok := extraConfig["claim_value"]; !ok || val != "some-value" {
			t.Errorf("Expected extra_config['claim_value'] to be 'some-value', got %v", val)
		}
		if _, ok := extraConfig["claim.value"]; ok {
			t.Errorf("Expected extra_config['claim.value'] to be removed")
		}
	})
}
