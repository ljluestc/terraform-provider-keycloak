# Fix: keycloak_custom_identity_provider_mapper not adding an attribute value

## Description
This PR fixes an issue where `attribute_value` and `claim_value` in the `extra_config` block of `keycloak_custom_identity_provider_mapper` were not being correctly passed to Keycloak. Keycloak expects these keys to be in dot notation (`attribute.value`, `claim.value`), but Terraform HCL maps usually prefer underscores or require quoting for dots.

The fix implements a transformation layer in `resource_keycloak_custom_identity_provider_mapper.go` that:
1.  Converts `attribute_value` -> `attribute.value` and `claim_value` -> `claim.value` before sending to Keycloak.
2.  Converts `attribute.value` -> `attribute_value` and `claim.value` -> `claim_value` when reading from Keycloak state.

This allows users to use standard Terraform snake_case convention in their configuration while ensuring Keycloak receives the expected dot notation.

## Issue
Fixes #806

## Type of Change
- [x] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Refactor (no functional changes, no api changes)
- [ ] Documentation
- [x] Test

## Checklist
- [x] I have performed a self-review of my own code
- [x] I have added tests that prove my fix is effective or that my feature works
- [x] New and existing unit tests pass locally with my changes
- [x] I have commented my code, particularly in hard-to-understand areas

## Verification
A new unit test `TestCustomIdentityProviderMapperKeyTransformation` has been added to `provider/resource_keycloak_custom_identity_provider_mapper_logic_test.go` to verify the transformation logic without requiring a running Keycloak instance.

To run the test:
```bash
SKIP_KEYCLOAK_SETUP=true go test -v ./provider -run TestCustomIdentityProviderMapperKeyTransformation
```
