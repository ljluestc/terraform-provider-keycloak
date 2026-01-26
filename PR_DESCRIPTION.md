Allow dotted keys in keycloak_custom_identity_provider_mapper extra_config

## Description
This PR addresses Issue #806 where users were unable to set attribute values for `hardcoded-attribute-idp-mapper` correctly. The user reported that providing `attribute_value` in `extra_config` did not result in the value being persisted in Keycloak.

Investigation revealed that this is a configuration mismatch. Keycloak's JSON API for the `hardcoded-attribute-idp-mapper` expects the key `attribute.value` (dotted), whereas the user was providing `attribute_value` (snake_case). The `keycloak_custom_identity_provider_mapper` resource passes `extra_config` keys directly to the API without transformation.

This PR:

*   **Verifies Support**: Adds a test case ensuring that `extra_config` keys containing dots (e.g., `attribute.value`) are correctly handled by the provider.
*   **Improves Documentation**: Updates the resource documentation to explicitly state that `extra_config` keys must match Keycloak's JSON property names exactly, including dots.

## Related Issue
Fixes #806

## Type of Change
- [x] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [x] Documentation update
- [x] Test addition (verification of existing functionality)

## Checklist
- [x] My code follows the style guidelines of this project
- [x] I have performed a self-review of my own code
- [x] I have commented my code, particularly in hard-to-understand areas
- [x] I have made corresponding changes to the documentation
- [x] I have added tests that prove my fix is effective or that my feature works
- [x] New and existing unit tests pass locally with my changes
