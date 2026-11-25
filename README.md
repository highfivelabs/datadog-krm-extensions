# datadog-krm-extensions

Kubernetes kinds for the Datadog configuration that Terraform's
`DataDog/datadog` provider manages and Datadog's operators do not. Datadog's official operator manages monitors, dashboards and SLOs (DatadogMonitor, DatadogDashboard, DatadogSLO, DatadogGenericResource) and deploys the Agent.

By the numbers: `DataDog/datadog` is the #17 Terraform provider by downloads
(630,210,450); version 4.30.0 has 178 resources. Operators
cover 5. This repository adds **166 kinds** standing in for
**171 more**. Each gap kind is listed below with its Terraform equivalent.

**Readiness: proven against a fake Datadog API only.** Every kind passes
create, update, drift repair and delete against a fake built from the same
descriptors (`TestEveryKindOnFake`). None has run against the real API yet;
treat each as untested there until it has.

## How it works

The kinds are generated, not hand-written. `hack/gen/vendor.yaml` pins
Datadog's OpenAPI spec; [krm-rest-runtime](https://github.com/highfivelabs/krm-rest-runtime)'s
generator finds every resource the API exposes in a create/read/update/delete
shape, keeps those that stand in for a Terraform resource in the gap list
(`hack/gen/terraform-gaps.json`), and writes:

- `types.yaml`: per kind, its API calls, path parameters and Terraform equivalents
- `config/crds/`: one CRD per kind, its spec the API's request body
- `zz_generated.report.md`: the full list

One generic manager (`cmd/manager`) runs every kind on krm-rest-runtime's
reconciler: it creates the resource, reads it back every 10 minutes and
puts back anything changed outside Kubernetes, and deletes it with the
object. Spec fields keep Datadog's own names, so the API reference
applies field for field. Path parameters are spec fields too.

Hand corrections to the automatic matching (a Terraform type matched to
the wrong API path, or not matched) are in `vendor.yaml`'s `overrides`.
Calls the vendor spec does not describe the way the generator needs them
are in hand-made OpenAPI supplements, `hack/gen/supplement/`, each citing
its sources (the pinned spec, Datadog's API docs, the Terraform provider's
calls) at the top.

## Install

```sh
kubectl kustomize config/default | KO_DOCKER_REPO=<registry> ko resolve -f - | kubectl apply --server-side -f -
```

Then, in each namespace that manages Datadog, a Secret with
`apiKey`, `appKey` and a Connection named `default`
(`config/samples/connection.yaml`). Objects use it unless
`spec.connectionRef.name` names another. Secret values (tokens, passwords,
secret values) go in a Kubernetes Secret, not the spec:
`spec.secretRefs: {<field>: {name, key}}` fills that body field.

To run only some API groups, install those groups' CRDs from `config/crds`
and pass `--groups` to the manager.

## Kinds

| Kind | Terraform equivalent | Shape | Read |
| --- | --- | --- | --- |
| `action.datadog.hfl.dev/v1alpha1` `ActionConnection` | `datadog_action_connection` | crud | `GET /api/v2/actions/connections/{connection_id}` |
| `action.datadog.hfl.dev/v1alpha1` `ActionExecutionPolicy` | `datadog_action_execution_policy` | crud | `GET /api/v2/actions/execution-policies/{policy_id}` |
| `agentless.datadog.hfl.dev/v1alpha1` `AgentlessScanningAwsScanOptions` | `datadog_agentless_scanning_aws_scan_options` | crud | `GET /api/v2/agentless_scanning/accounts/aws/{account_id}` |
| `agentless.datadog.hfl.dev/v1alpha1` `AgentlessScanningAzureScanOptions` | `datadog_agentless_scanning_azure_scan_options` | crud | `GET /api/v2/agentless_scanning/accounts/azure/{subscription_id}` |
| `agentless.datadog.hfl.dev/v1alpha1` `AgentlessScanningGcpScanOptions` | `datadog_agentless_scanning_gcp_scan_options` | crud | `GET /api/v2/agentless_scanning/accounts/gcp/{project_id}` |
| `api.datadog.hfl.dev/v1alpha1` `ApiKey` | `datadog_api_key` | crud | `GET /api/v1/api_key/{key}` |
| `apm.datadog.hfl.dev/v1alpha1` `ApmRetentionFilter` | `datadog_apm_retention_filter` | crud | `GET /api/v2/apm/config/retention-filters/{filter_id}` |
| `apm.datadog.hfl.dev/v1alpha1` `ApmRetentionFilterOrder` | `datadog_apm_retention_filter_order` | settings | `GET /api/v2/apm/config/retention-filters` |
| `app.datadog.hfl.dev/v1alpha1` `AppBuilderApp` | `datadog_app_builder_app` | crud | `GET /api/v2/app-builder/apps/{app_id}` |
| `app.datadog.hfl.dev/v1alpha1` `AppKeyRegistration` | `datadog_app_key_registration` | upsert | `GET /api/v2/actions/app_key_registrations/{app_key_id}` |
| `application.datadog.hfl.dev/v1alpha1` `ApplicationKey` | `datadog_application_key` | crud | `GET /api/v2/current_user/application_keys/{app_key_id}` |
| `appsec.datadog.hfl.dev/v1alpha1` `AppsecWafCustomRule` | `datadog_appsec_waf_custom_rule` | crud | `GET /api/v2/remote_config/products/asm/waf/custom_rules/{custom_rule_id}` |
| `appsec.datadog.hfl.dev/v1alpha1` `AppsecWafExclusionFilter` | `datadog_appsec_waf_exclusion_filter` | crud | `GET /api/v2/remote_config/products/asm/waf/exclusion_filters/{exclusion_filter_id}` |
| `authn.datadog.hfl.dev/v1alpha1` `AuthnMapping` | `datadog_authn_mapping` | crud | `GET /api/v2/authn_mappings/{authn_mapping_id}` |
| `aws.datadog.hfl.dev/v1alpha1` `AwsCurConfig` | `datadog_aws_cur_config` | crud | `GET /api/v2/cost/aws_cur_config/{cloud_account_id}` |
| `aws.datadog.hfl.dev/v1alpha1` `AwsWifIdentityMapping` | `datadog_aws_wif_identity_mapping` | crud | `GET /api/v2/wif/aws/persona_mapping/{persona_mapping_id}` |
| `aws.datadog.hfl.dev/v1alpha1` `AwsWifIntakeMapping` | `datadog_aws_wif_intake_mapping` | crud | `GET /api/v2/wif/aws/intake_mapping/{intake_mapping_id}` |
| `azure.datadog.hfl.dev/v1alpha1` `AzureUcConfig` | `datadog_azure_uc_config` | crud | `GET /api/v2/cost/azure_uc_config/{cloud_account_id}` |
| `child.datadog.hfl.dev/v1alpha1` `ChildOrganization` | `datadog_child_organization` | crud | `GET /api/v2/org` |
| `cloud.datadog.hfl.dev/v1alpha1` `CloudInventorySyncConfig` | `datadog_cloud_inventory_sync_config` | crud | `GET /api/v2/cloudinventoryservice/syncconfigs/{id}` |
| `compliance.datadog.hfl.dev/v1alpha1` `ComplianceCustomFramework` | `datadog_compliance_custom_framework` | crud | `GET /api/v2/cloud_security_management/custom_frameworks/{handle}/{version}` |
| `compliance.datadog.hfl.dev/v1alpha1` `ComplianceResourceEvaluationFilter` | `datadog_compliance_resource_evaluation_filter` | settings | `GET /api/v2/cloud_security_management/resource_filters` |
| `cost.datadog.hfl.dev/v1alpha1` `CostBudget` | `datadog_cost_budget` | crud | `GET /api/v2/cost/budget/{budget_id}` |
| `cost.datadog.hfl.dev/v1alpha1` `CostCustomForecast` | `datadog_cost_custom_forecast` | upsert | `GET /api/v2/cost/budget/{budget_id}/custom-forecast` |
| `csm.datadog.hfl.dev/v1alpha1` `CsmThreatsAgentRule` | `datadog_cloud_workload_security_agent_rule`, `datadog_csm_threats_agent_rule` | crud | `GET /api/v2/remote_config/products/cws/agent_rules/{agent_rule_id}` |
| `csm.datadog.hfl.dev/v1alpha1` `CsmThreatsPolicy` | `datadog_csm_threats_policy` | crud | `GET /api/v2/remote_config/products/cws/policy/{policy_id}` |
| `custom.datadog.hfl.dev/v1alpha1` `CustomAllocationRule` | `datadog_custom_allocation_rule` | crud | `GET /api/v2/cost/arbitrary_rule/{rule_id}` |
| `custom.datadog.hfl.dev/v1alpha1` `CustomAllocationRuleOrder` | `datadog_custom_allocation_rules` | settings | `GET /api/v2/cost/arbitrary_rule` |
| `dashboard.datadog.hfl.dev/v1alpha1` `DashboardList` | `datadog_dashboard_list` | crud | `GET /api/v1/dashboard/lists/manual/{list_id}` |
| `dashboard.datadog.hfl.dev/v1alpha1` `DashboardV2` | `datadog_dashboard_json`, `datadog_dashboard_v2` | crud | `GET /api/v1/dashboard/{dashboard_id}` |
| `dataset.datadog.hfl.dev/v1alpha1` `Dataset` | `datadog_dataset` | crud | `GET /api/v2/datasets/{dataset_id}` |
| `datastore.datadog.hfl.dev/v1alpha1` `Datastore` | `datadog_datastore` | crud | `GET /api/v2/actions-datastores/{datastore_id}` |
| `datastore.datadog.hfl.dev/v1alpha1` `DatastoreItem` | `datadog_datastore_item` | upsert | `GET /api/v2/actions-datastores/{datastore_id}/items` |
| `deployment.datadog.hfl.dev/v1alpha1` `DeploymentGate` | `datadog_deployment_gate` | crud | `GET /api/v2/deployment_gates/{id}` |
| `derived.datadog.hfl.dev/v1alpha1` `DerivedMetric` | `datadog_derived_metric` | crud | `GET /api/v2/metrics/derived-metrics/{metric_name}` |
| `domain.datadog.hfl.dev/v1alpha1` `DomainAllowlist` | `datadog_domain_allowlist` | settings | `GET /api/v2/domain_allowlist` |
| `downtime.datadog.hfl.dev/v1alpha1` `Downtime` | `datadog_downtime` | crud | `GET /api/v1/downtime/{downtime_id}` |
| `fleet.datadog.hfl.dev/v1alpha1` `FleetSchedule` | `datadog_fleet_schedule` | crud | `GET /api/v2/fleet/schedules/{id}` |
| `gcp.datadog.hfl.dev/v1alpha1` `GcpUcConfig` | `datadog_gcp_uc_config` | crud | `GET /api/v2/cost/gcp_uc_config/{cloud_account_id}` |
| `github.datadog.hfl.dev/v1alpha1` `GithubOidcIdentityMapping` | `datadog_github_oidc_identity_mapping` | crud | `GET /api/v2/wif/github/persona_mapping/{persona_mapping_id}` |
| `github.datadog.hfl.dev/v1alpha1` `GithubOidcIntakeMapping` | `datadog_github_oidc_intake_mapping` | crud | `GET /api/v2/wif/github/intake_mapping/{intake_mapping_id}` |
| `governance.datadog.hfl.dev/v1alpha1` `GovernanceControl` | `datadog_governance_control` | settings | `GET /api/v2/governance/control/{detection_type}` |
| `incident.datadog.hfl.dev/v1alpha1` `IncidentNotificationRule` | `datadog_incident_notification_rule` | crud | `GET /api/v2/incidents/config/notification-rules/{id}` |
| `incident.datadog.hfl.dev/v1alpha1` `IncidentNotificationTemplate` | `datadog_incident_notification_template` | crud | `GET /api/v2/incidents/config/notification-templates/{id}` |
| `incident.datadog.hfl.dev/v1alpha1` `IncidentPostmortemTemplate` | `datadog_incident_postmortem_template` | crud | `GET /api/v2/incidents/config/postmortem-templates/{template_id}` |
| `incident.datadog.hfl.dev/v1alpha1` `IncidentType` | `datadog_incident_type` | crud | `GET /api/v2/incidents/config/types/{incident_type_id}` |
| `incident.datadog.hfl.dev/v1alpha1` `IncidentUserDefinedField` | `datadog_incident_user_defined_field` | crud | `GET /api/v2/incidents/config/user-defined-fields/{field_id}` |
| `incident.datadog.hfl.dev/v1alpha1` `IncidentUserDefinedRole` | `datadog_incident_user_defined_role` | crud | `GET /api/v2/incidents/config/user-defined-roles/{role_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationAwsAccount` | `datadog_integration_aws_account` | crud | `GET /api/v2/integration/aws/accounts/{aws_account_config_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationAwsAccountCcmConfig` | `datadog_integration_aws_account_ccm_config` | upsert | `GET /api/v2/integration/aws/accounts/{aws_account_config_id}/ccm_config` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationAwsEventBridge` | `datadog_integration_aws_event_bridge` | crud | `GET /api/v2/integration/aws/event_bridge` (finds `name`) |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationAzure` | `datadog_integration_azure` | upsert | `GET /api/v1/integration/azure` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationCloudflareAccount` | `datadog_integration_cloudflare_account` | crud | `GET /api/v2/integrations/cloudflare/accounts/{account_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationConfluentAccount` | `datadog_integration_confluent_account` | crud | `GET /api/v2/integrations/confluent-cloud/accounts/{account_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationConfluentResource` | `datadog_integration_confluent_resource` | crud | `GET /api/v2/integrations/confluent-cloud/accounts/{account_id}/resources/{resource_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationDatabricksAccount` | `datadog_integration_databricks_account` | crud | `GET /api/v2/integration-interfaces/databricks/accounts/{account_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationElasticCloudAccount` | `datadog_integration_elastic_cloud_account` | crud | `GET /api/v2/integration-interfaces/elastic-cloud/accounts/{account_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationFastlyAccount` | `datadog_integration_fastly_account` | crud | `GET /api/v2/integrations/fastly/accounts/{account_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationFastlyService` | `datadog_integration_fastly_service` | crud | `GET /api/v2/integrations/fastly/accounts/{account_id}/services/{service_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationGcp` | `datadog_integration_gcp` | upsert | `GET /api/v1/integration/gcp` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationGcpSts` | `datadog_integration_gcp_sts` | crud | `GET /api/v2/integration/gcp/accounts` (finds `id`) |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationMsTeamsTenantBasedHandle` | `datadog_integration_ms_teams_tenant_based_handle` | crud | `GET /api/v2/integration/ms-teams/configuration/tenant-based-handles/{handle_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationMsTeamsWorkflowsWebhookHandle` | `datadog_integration_ms_teams_workflows_webhook_handle` | crud | `GET /api/v2/integration/ms-teams/configuration/workflows-webhook-handles/{handle_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationOpsgenieServiceObject` | `datadog_integration_opsgenie_service_object` | crud | `GET /api/v2/integration/opsgenie/services/{integration_service_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationPagerduty` | `datadog_integration_pagerduty_service_object` | crud | `GET /api/v1/integration/pagerduty/configuration/services/{service_name}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationPagerdutySettings` | `datadog_integration_pagerduty` | settings | `GET /api/v1/integration/pagerduty` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationSlackChannel` | `datadog_integration_slack_channel` | crud | `GET /api/v1/integration/slack/configuration/accounts/{account_name}/channels/{channel_name}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationSnowflakeAccount` | `datadog_integration_snowflake_account` | crud | `GET /api/v2/integration-interfaces/snowflake/accounts/{account_id}` |
| `integration.datadog.hfl.dev/v1alpha1` `IntegrationTwilioAccount` | `datadog_integration_twilio_account` | crud | `GET /api/v2/integration-interfaces/twilio/accounts/{account_id}` |
| `ip.datadog.hfl.dev/v1alpha1` `IpAllowlist` | `datadog_ip_allowlist` | settings | `GET /api/v2/ip_allowlist` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsArchive` | `datadog_logs_archive` | crud | `GET /api/v2/logs/config/archives/{archive_id}` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsArchiveOrder` | `datadog_logs_archive_order` | settings | `GET /api/v2/logs/config/archive-order` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsCustomDestination` | `datadog_logs_custom_destination` | crud | `GET /api/v2/logs/config/custom-destinations/{custom_destination_id}` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsCustomPipeline` | `datadog_logs_custom_pipeline` | crud | `GET /api/v1/logs/config/pipelines/{pipeline_id}` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsIndex` | `datadog_logs_index` | crud | `GET /api/v1/logs/config/indexes/{name}` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsIndexOrder` | `datadog_logs_index_order` | settings | `GET /api/v1/logs/config/index-order` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsIntegrationPipeline` | `datadog_logs_integration_pipeline` | settings | `GET /api/v1/logs/config/pipelines/{pipeline_id}` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsMetric` | `datadog_logs_metric` | crud | `GET /api/v2/logs/config/metrics/{metric_id}` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsPipelineOrder` | `datadog_logs_pipeline_order` | settings | `GET /api/v1/logs/config/pipeline-order` |
| `logs.datadog.hfl.dev/v1alpha1` `LogsRestrictionQuery` | `datadog_logs_restriction_query` | crud | `GET /api/v2/logs/config/restriction_queries/{restriction_query_id}` |
| `metric.datadog.hfl.dev/v1alpha1` `MetricMetadata` | `datadog_metric_metadata` | settings | `GET /api/v1/metrics/{metric_name}` |
| `metric.datadog.hfl.dev/v1alpha1` `MetricTagConfiguration` | `datadog_metric_tag_configuration` | upsert | `GET /api/v2/metrics/{metric_name}/tags` |
| `monitor.datadog.hfl.dev/v1alpha1` `Monitor` | `datadog_monitor_json` | crud | `GET /api/v1/monitor/{monitor_id}` |
| `monitor.datadog.hfl.dev/v1alpha1` `MonitorConfigPolicy` | `datadog_monitor_config_policy` | crud | `GET /api/v2/monitor/policy/{policy_id}` |
| `monitor.datadog.hfl.dev/v1alpha1` `MonitorNotificationRule` | `datadog_monitor_notification_rule` | crud | `GET /api/v2/monitor/notification_rule/{rule_id}` |
| `observability.datadog.hfl.dev/v1alpha1` `ObservabilityPipeline` | `datadog_observability_pipeline` | crud | `GET /api/v2/obs-pipelines/pipelines/{pipeline_id}` |
| `on.datadog.hfl.dev/v1alpha1` `OnCallEscalationPolicy` | `datadog_on_call_escalation_policy` | crud | `GET /api/v2/on-call/escalation-policies/{policy_id}` |
| `on.datadog.hfl.dev/v1alpha1` `OnCallSchedule` | `datadog_on_call_schedule` | crud | `GET /api/v2/on-call/schedules/{schedule_id}` |
| `on.datadog.hfl.dev/v1alpha1` `OnCallTeamRoutingRules` | `datadog_on_call_team_routing_rules` | settings | `GET /api/v2/on-call/teams/{team_id}/routing-rules` |
| `on.datadog.hfl.dev/v1alpha1` `OnCallUserNotificationChannel` | `datadog_on_call_user_notification_channel` | crud | `GET /api/v2/on-call/users/{user_id}/notification-channels/{channel_id}` |
| `on.datadog.hfl.dev/v1alpha1` `OnCallUserNotificationRule` | `datadog_on_call_user_notification_rule` | crud | `GET /api/v2/on-call/users/{user_id}/notification-rules/{rule_id}` |
| `openapi.datadog.hfl.dev/v1alpha1` `OpenapiApi` | `datadog_openapi_api` | settings | `GET /api/v2/apicatalog/api/{id}/openapi` |
| `org.datadog.hfl.dev/v1alpha1` `OrgConnection` | `datadog_org_connection` | crud | `GET /api/v2/org_connections` (finds `id`) |
| `org.datadog.hfl.dev/v1alpha1` `OrgGroup` | `datadog_org_group` | crud | `GET /api/v2/org_groups/{org_group_id}` |
| `org.datadog.hfl.dev/v1alpha1` `OrgGroupMembership` | `datadog_org_group_membership` | crud | `GET /api/v2/org_group_memberships/{org_group_membership_id}` |
| `org.datadog.hfl.dev/v1alpha1` `OrgGroupPolicy` | `datadog_org_group_policy` | crud | `GET /api/v2/org_group_policies/{org_group_policy_id}` |
| `org.datadog.hfl.dev/v1alpha1` `OrgGroupPolicyOverride` | `datadog_org_group_policy_override` | crud | `GET /api/v2/org_group_policy_overrides/{org_group_policy_override_id}` |
| `organization.datadog.hfl.dev/v1alpha1` `OrganizationSettings` | `datadog_organization_settings` | settings | `GET /api/v1/org/{public_id}` |
| `powerpack.datadog.hfl.dev/v1alpha1` `Powerpack` | `datadog_powerpack`, `datadog_powerpack_v2` | crud | `GET /api/v2/powerpacks/{powerpack_id}` |
| `reference.datadog.hfl.dev/v1alpha1` `ReferenceTable` | `datadog_reference_table` | crud | `GET /api/v2/reference-tables/tables/{id}` |
| `restriction.datadog.hfl.dev/v1alpha1` `RestrictionPolicy` | `datadog_restriction_policy` | upsert | `GET /api/v2/restriction_policy/{resource_id}` |
| `role.datadog.hfl.dev/v1alpha1` `Role` | `datadog_role` | crud | `GET /api/v2/roles/{role_id}` |
| `rum.datadog.hfl.dev/v1alpha1` `RumApplication` | `datadog_rum_application` | crud | `GET /api/v2/rum/applications/{id}` |
| `rum.datadog.hfl.dev/v1alpha1` `RumExclusionFilter` | `datadog_rum_exclusion_filter` | crud | `GET /api/v2/rum/applications/{app_id}/retention_filters/exclusion/{ef_id}` |
| `rum.datadog.hfl.dev/v1alpha1` `RumMetric` | `datadog_rum_metric` | crud | `GET /api/v2/rum/config/metrics/{metric_id}` |
| `rum.datadog.hfl.dev/v1alpha1` `RumOperation` | `datadog_rum_operation` | crud | `GET /api/v2/rum/operations/{rum_operation_id}` |
| `rum.datadog.hfl.dev/v1alpha1` `RumOperationJourneyLink` | `datadog_rum_operation_journey_link` | crud | `GET /api/v2/rum/operations/strong_links` |
| `rum.datadog.hfl.dev/v1alpha1` `RumRetentionFilter` | `datadog_rum_retention_filter` | crud | `GET /api/v2/rum/applications/{app_id}/retention_filters/{rf_id}` |
| `rum.datadog.hfl.dev/v1alpha1` `RumRetentionFiltersOrder` | `datadog_rum_retention_filters_order` | settings | `GET /api/v2/rum/applications/{app_id}/retention_filters` |
| `rum.datadog.hfl.dev/v1alpha1` `RumRetentionQuota` | `datadog_rum_retention_quota` | upsert | `GET /api/v2/rum/config/retention-quota/{scope_type}/{scope_id}` |
| `secure.datadog.hfl.dev/v1alpha1` `SecureEmbedDashboard` | `datadog_secure_embed_dashboard` | crud | `GET /api/v2/dashboard/{dashboard_id}/shared/secure-embed/{token}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsDefaultInboxRule` | `datadog_security_findings_default_inbox_rule` | settings | `GET /api/v2/security/findings/automation/default_inbox_rules/{rule_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsDueDateRule` | `datadog_security_findings_due_date_rule` | crud | `GET /api/v2/security/findings/automation/due_date_rules/{rule_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsDueDateRulesOrder` | `datadog_security_findings_due_date_rules_order` | settings | `GET /api/v2/security/findings/automation/due_date_rules` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsInboxRule` | `datadog_security_findings_inbox_rule` | crud | `GET /api/v2/security/findings/automation/inbox_rules/{rule_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsInboxRulesOrder` | `datadog_security_findings_inbox_rules_order` | settings | `GET /api/v2/security/findings/automation/inbox_rules` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsMuteRule` | `datadog_security_findings_mute_rule` | crud | `GET /api/v2/security/findings/automation/mute_rules/{rule_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsMuteRulesOrder` | `datadog_security_findings_mute_rules_order` | settings | `GET /api/v2/security/findings/automation/mute_rules` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsSeverityModifierRule` | `datadog_security_findings_severity_modifier_rule` | crud | `GET /api/v2/security/findings/automation/severity_modifier_rules/{rule_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsSeverityModifierRulesOrder` | `datadog_security_findings_severity_modifier_rules_order` | settings | `GET /api/v2/security/findings/automation/severity_modifier_rules` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsTicketCreationRule` | `datadog_security_findings_ticket_creation_rule` | crud | `GET /api/v2/security/findings/automation/ticket_creation_rules/{rule_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityFindingsTicketCreationRulesOrder` | `datadog_security_findings_ticket_creation_rules_order` | settings | `GET /api/v2/security/findings/automation/ticket_creation_rules` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityMonitoringCriticalAsset` | `datadog_security_monitoring_critical_asset` | crud | `GET /api/v2/security_monitoring/configuration/critical_assets/{critical_asset_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityMonitoringDefaultRule` | `datadog_security_monitoring_default_rule` | settings | `GET /api/v2/security_monitoring/rules/{rule_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityMonitoringFilter` | `datadog_security_monitoring_filter` | crud | `GET /api/v2/security_monitoring/configuration/security_filters/{security_filter_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityMonitoringRule` | `datadog_cloud_configuration_rule`, `datadog_security_monitoring_rule`, `datadog_security_monitoring_rule_json` | crud | `GET /api/v2/security_monitoring/rules/{rule_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityMonitoringSuppression` | `datadog_security_monitoring_suppression` | crud | `GET /api/v2/security_monitoring/configuration/suppressions/{suppression_id}` |
| `security.datadog.hfl.dev/v1alpha1` `SecurityNotificationRule` | `datadog_security_notification_rule` | crud | `GET /api/v2/security/signals/notification_rules/{id}` |
| `sensitive.datadog.hfl.dev/v1alpha1` `SensitiveDataScannerGroup` | `datadog_sensitive_data_scanner_group` | crud | `GET /api/v2/sensitive-data-scanner/config` (finds `id` in `included`) |
| `sensitive.datadog.hfl.dev/v1alpha1` `SensitiveDataScannerGroupOrder` | `datadog_sensitive_data_scanner_group_order` | settings | `GET /api/v2/sensitive-data-scanner/config` |
| `sensitive.datadog.hfl.dev/v1alpha1` `SensitiveDataScannerRule` | `datadog_sensitive_data_scanner_rule` | crud | `GET /api/v2/sensitive-data-scanner/config` (finds `id` in `included`) |
| `service.datadog.hfl.dev/v1alpha1` `ServiceAccessToken` | `datadog_service_access_token` | crud | `GET /api/v2/service_accounts/{service_account_id}/access_tokens/{token_id}` |
| `service.datadog.hfl.dev/v1alpha1` `ServiceAccount` | `datadog_service_account` | crud | `GET /api/v2/users/{service_account_id}` |
| `service.datadog.hfl.dev/v1alpha1` `ServiceAccountApplicationKey` | `datadog_service_account_application_key` | crud | `GET /api/v2/service_accounts/{service_account_id}/application_keys/{app_key_id}` |
| `service.datadog.hfl.dev/v1alpha1` `ServiceDefinitionYaml` | `datadog_service_definition_yaml` | upsert | `GET /api/v2/services/definitions/{service_name}` |
| `slo.datadog.hfl.dev/v1alpha1` `SloCorrection` | `datadog_slo_correction` | crud | `GET /api/v1/slo/correction/{slo_correction_id}` |
| `software.datadog.hfl.dev/v1alpha1` `SoftwareCatalog` | `datadog_software_catalog` | crud | `GET /api/v2/catalog/entity` |
| `spans.datadog.hfl.dev/v1alpha1` `SpansMetric` | `datadog_spans_metric` | crud | `GET /api/v2/apm/config/metrics/{metric_id}` |
| `status.datadog.hfl.dev/v1alpha1` `StatusPage` | `datadog_status_page` | crud | `GET /api/v2/statuspages/{page_id}` |
| `status.datadog.hfl.dev/v1alpha1` `StatusPageComponent` | `datadog_status_page_component` | crud | `GET /api/v2/statuspages/{page_id}/components/{component_id}` |
| `status.datadog.hfl.dev/v1alpha1` `StatusPageDegradationTemplate` | `datadog_status_page_degradation_template` | crud | `GET /api/v2/statuspages/{page_id}/degradation_templates/{template_id}` |
| `status.datadog.hfl.dev/v1alpha1` `StatusPageMaintenanceTemplate` | `datadog_status_page_maintenance_template` | crud | `GET /api/v2/statuspages/{page_id}/maintenance_templates/{template_id}` |
| `synthetics.datadog.hfl.dev/v1alpha1` `SyntheticsConcurrencyCap` | `datadog_synthetics_concurrency_cap` | settings | `GET /api/v2/synthetics/settings/on_demand_concurrency_cap` |
| `synthetics.datadog.hfl.dev/v1alpha1` `SyntheticsGlobalVariable` | `datadog_synthetics_global_variable` | crud | `GET /api/v1/synthetics/variables/{variable_id}` |
| `synthetics.datadog.hfl.dev/v1alpha1` `SyntheticsPrivateLocation` | `datadog_synthetics_private_location` | crud | `GET /api/v1/synthetics/private-locations/{location_id}` |
| `synthetics.datadog.hfl.dev/v1alpha1` `SyntheticsSuite` | `datadog_synthetics_suite` | crud | `GET /api/v2/synthetics/suites/{public_id}` |
| `tag.datadog.hfl.dev/v1alpha1` `TagIndexingRule` | `datadog_tag_indexing_rule` | crud | `GET /api/v2/metrics/tag-indexing-rules/{id}` |
| `tag.datadog.hfl.dev/v1alpha1` `TagIndexingRuleExemption` | `datadog_tag_indexing_rule_exemption` | upsert | `GET /api/v2/metrics/{metric_name}/tag-indexing-rule-exemptions` |
| `tag.datadog.hfl.dev/v1alpha1` `TagIndexingRuleOrder` | `datadog_tag_indexing_rule_order` | settings | `GET /api/v2/metrics/tag-indexing-rules` |
| `tag.datadog.hfl.dev/v1alpha1` `TagPipelineRuleset` | `datadog_tag_pipeline_ruleset` | crud | `GET /api/v2/tags/enrichment/{ruleset_id}` |
| `tag.datadog.hfl.dev/v1alpha1` `TagPipelineRulesetOrder` | `datadog_tag_pipeline_rulesets` | settings | `GET /api/v2/tags/enrichment` |
| `tag.datadog.hfl.dev/v1alpha1` `TagRule` | `datadog_tag_rule` | crud | `GET /api/v2/governance/tag_rules/{rule_id}` |
| `team.datadog.hfl.dev/v1alpha1` `Team` | `datadog_team` | crud | `GET /api/v2/team/{team_id}` |
| `team.datadog.hfl.dev/v1alpha1` `TeamConnection` | `datadog_team_connection` | crud | `GET /api/v2/team/connections` (finds `id`) |
| `team.datadog.hfl.dev/v1alpha1` `TeamHierarchyLinks` | `datadog_team_hierarchy_links` | crud | `GET /api/v2/team-hierarchy-links/{link_id}` |
| `team.datadog.hfl.dev/v1alpha1` `TeamLink` | `datadog_team_link` | crud | `GET /api/v2/team/{team_id}/links/{link_id}` |
| `team.datadog.hfl.dev/v1alpha1` `TeamMembership` | `datadog_team_membership` | crud | `GET /api/v2/team/{team_id}/memberships` (finds `relationships.user.data.id`) |
| `team.datadog.hfl.dev/v1alpha1` `TeamNotificationRule` | `datadog_team_notification_rule` | crud | `GET /api/v2/team/{team_id}/notification-rules/{rule_id}` |
| `team.datadog.hfl.dev/v1alpha1` `TeamPermissionSetting` | `datadog_team_permission_setting` | settings | `GET /api/v2/team/{team_id}/permission-settings` (finds `attributes.action`) |
| `team.datadog.hfl.dev/v1alpha1` `TeamSync` | `datadog_team_sync` | upsert | `GET /api/v2/team/sync` |
| `user.datadog.hfl.dev/v1alpha1` `User` | `datadog_user` | crud | `GET /api/v2/users/{user_id}` |
| `user.datadog.hfl.dev/v1alpha1` `UserRole` | `datadog_user_role` | assoc | `GET /api/v2/roles/{role_id}/users` |
| `webhook.datadog.hfl.dev/v1alpha1` `Webhook` | `datadog_webhook` | crud | `GET /api/v1/integration/webhooks/configuration/webhooks/{webhook_name}` |
| `webhook.datadog.hfl.dev/v1alpha1` `WebhookCustomVariable` | `datadog_webhook_custom_variable` | crud | `GET /api/v1/integration/webhooks/configuration/custom-variables/{custom_variable_name}` |
| `webhook.datadog.hfl.dev/v1alpha1` `WebhookOauth2ClientCredentials` | `datadog_webhook_oauth2_client_credentials` | crud | `GET /api/v2/integration/webhooks/configuration/auth-method/oauth2-client-credentials/{auth_method_id}` |
| `workflow.datadog.hfl.dev/v1alpha1` `WorkflowAutomation` | `datadog_workflow_automation` | crud | `GET /api/v2/workflows/{workflow_id}` |


Kinds read from a list (the API has no GET for one item) find the item in
the list's first page, as the "finds" column says: `TeamMembership` asks
for 100 members a page, so a team with more members than that may not
show the one it looks for. `IntegrationGcpSts` needs the organization's
Datadog delegate service account, which the Terraform provider makes
before each create (`POST /api/v2/integration/gcp/sts_delegate`, once per
organization, also made by setting up the GCP integration in Datadog).

`SensitiveDataScannerGroup` and `SensitiveDataScannerRule` send
`{"meta": {}}` beside `data` on every write, and as their DELETE's body.
`CostBudget` is created and replaced by `PUT /api/v2/cost/budget`, the
budget's ID in `data.id` once it has one. Deleting a `SyntheticsSuite` is a
`POST /api/v2/synthetics/suites/bulk-delete` naming it. `IntegrationAzure`
and `IntegrationGcp` (the v1 integrations) have no item path: each is named
by `spec.tenant_name` or `spec.project_id` on the collection, made by a
POST, found in the GET's list, changed by a PUT and removed by a DELETE
whose body names it.

## Not covered

<!-- BEGIN generated: not covered (krmgen.py, from vendor.yaml hardBlocks and draftNotes) -->
2 of the 173 Terraform gap types have no kind: 2 are not built (`hardBlocks` in `hack/gen/vendor.yaml`).

| Terraform | Status | Why |
| --- | --- | --- |
| `datadog_integration_aws_external_id` | not built | Not a resource (nothing to read back or delete): generates a new external ID; an action, not an object. |
| `datadog_saml_idp_metadata` | not built | Not a resource (nothing to read back or delete): uploads IdP metadata; an action. |
<!-- END generated: not covered -->

## Kinds with a limit worth knowing

- `ChildOrganization` makes a child organization (`POST /api/v1/org`) and
  reads it among the caller's managed organizations (`GET /api/v2/org`,
  found in `included` by `attributes.public_id`), which the parent's keys
  answer; the provider never reads one back. Nothing deletes an
  organization (Datadog support does): deleting the object leaves it. Its
  name cannot be changed (a new name is a new organization). The API and
  application keys the create answers (`api_key.key`,
  `application_key.hash`) go to `spec.writeSecretRef`, answered only then.
- `DatastoreItem` is one item of a datastore, named by `item_key`; each call
  builds its own body from the spec: added with `POST .../items/bulk`
  (`values: [value]`, failing on a conflict), changed with `PATCH .../items`
  (`item_changes.ops_set: value`), removed with `DELETE .../items`, read
  with `?item_key=`. Its `value` holds the datastore's primary key column.

- `LogsIntegrationPipeline`, `SecurityMonitoringDefaultRule` and
  `SecurityFindingsDefaultInboxRule` are objects Datadog makes: each names
  one by its ID (as the provider must import it) and sets what may change
  (whether it is on; a default rule's custom name, message and filters).
  Deleting the object leaves it. `SecurityMonitoringDefaultRule` does not
  set the provider's case notifications, query extensions, options or
  custom tags. `SecurityFindingsDefaultInboxRule`'s `enabled` is the string
  `"true"` or `"false"`, the call (`.../enable`, `.../disable`) it makes.
- `TeamSync` is one configuration per source; deleting the object pauses
  it (`frequency: paused`), as the provider does.
- `CloudInventorySyncConfig` is read back in another form than written, so
  its settings are sent when the spec changes, not compared.
- `SoftwareCatalog` and `ServiceDefinitionYaml` take the entity or
  definition as fields (JSON), where the provider takes YAML text; a
  software catalog entity needs `metadata.namespace` (Datadog's default is
  `default`).

- The `*Order` kinds (and `CustomAllocationRuleOrder`, `TagPipelineRulesetOrder`
  for `datadog_custom_allocation_rules` and `datadog_tag_pipeline_rulesets`)
  hold the order of other objects: their spec lists the objects' IDs, read
  back from the objects' list and written whole by the reorder call.
  Deleting one leaves the order as it is, as the provider does. Most take
  every object's ID, as the API does; `CustomAllocationRuleOrder` and
  `TagPipelineRulesetOrder` take those they order first and keep the others
  after them, as the provider does by default (its
  `override_ui_defined_resources`, which deletes the others, is not covered).
- `UserRole` is a link: its spec is the role (`roleID`) and the user (`id`,
  `type: users`); it is looked for in the role's first 100 users.
- `RumOperationJourneyLink` sets the link's `status` as the spec says (the
  provider always confirms it); a link that exists already as a draft is not
  taken over.
- `ComplianceCustomFramework` is named by its `attributes.handle` and
  `attributes.version` (the path `.../{handle}/{version}`): to change
  either, delete the object and make a new one.

## For agents (MCP)

`cmd/mcp` serves this repository's kinds over the Model Context Protocol
(stdio), from `types.yaml`, `config/crds` and the live tests' records.
`.mcp.json` registers it, so an agent session opened in this repository
has it without setup. Elsewhere, run it from the repository root:

```sh
go run ./cmd/mcp
```

Its tools: `list_kinds`, `find_kinds_for_terraform`, `describe_kind` (every
field, the REST calls, the live record), `example` (the live-tested spec
where there is one), `validate_manifest` (the CRD schema, as the API server
checks it) and `plan_on_fake` (the controller run against the fake API,
with every REST call it makes). Nothing it does reaches the real API.
krm-rest-runtime's README ("For agents") has the details.

## Regenerate

```sh
hack/gen/generate.sh   # needs python3 + pyyaml, and krm-rest-runtime beside this repo
go test ./...
```
