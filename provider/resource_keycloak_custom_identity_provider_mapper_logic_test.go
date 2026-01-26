package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestCustomIdentityProviderMapperExtraConfigPassthrough(t *testing.T) {
	resource := resourceKeycloakCustomIdentityProviderMapper()
	resourceSchema := resource.Schema

	// Test Case: Verify dotted keys are preserved (no transformation)
	t.Run("dotted keys preservation", func(t *testing.T) {
		rawData := map[string]interface{}{
			"realm":                    "my-realm",
			"name":                     "my-mapper",
			"identity_provider_alias":  "my-idp",
			"identity_provider_mapper": "hardcoded-attribute-idp-mapper",
			"extra_config": map[string]interface{}{
				"attribute":       "something",
				"attribute.value": "some-value", // Keycloak expects this
			},
		}

		data := schema.TestResourceDataRaw(t, resourceSchema, rawData)

		mapper := getCustomIdentityProviderMapperFromData(data)

		// Verify "attribute.value" is present and correct
		if val, ok := mapper.Config.ExtraConfig["attribute.value"]; !ok || val != "some-value" {
			t.Errorf("Expected extra_config['attribute.value'] to be 'some-value', got %v", val)
		}
		// Verify no magical "attribute_value" created
		if _, ok := mapper.Config.ExtraConfig["attribute_value"]; ok {
			t.Errorf("Expected extra_config['attribute_value'] NOT to be present (no transformation)")
		}

		// Verify reverse (read from Keycloak)
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

		dataForSet := schema.TestResourceDataRaw(t, resourceSchema, rawData)
		setCustomIdentityProviderMapperData(dataForSet, returnedMapper)

		extraConfig := dataForSet.Get("extra_config").(map[string]interface{})
		
		// Verify "attribute.value" is preserved in state
		if val, ok := extraConfig["attribute.value"]; !ok || val != "some-value" {
			t.Errorf("Expected extra_config['attribute.value'] to be 'some-value', got %v", val)
		}
	})
}
