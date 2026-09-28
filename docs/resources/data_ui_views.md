# Resource: splunk_data_ui_views
Create and manage splunk dashboards/views.
## Example Usage
```
resource "splunk_data_ui_views" "dashboard" {
  name     = "Terraform_Sample_Dashboard"
  eai_data = "<dashboard version=\"1.1\"><label>Terraform</label><description>Terraform operations</description><row><panel><chart><search><query>index=_internal sourcetype=splunkd_access useragent=\"splunk-simple-go-client\" | timechart fixedrange=f values(status) by uri_path</query><earliest>-24h@h</earliest><latest>now</latest><sampleRatio>1</sampleRatio></search><option name=\"charting.axisLabelsX.majorLabelStyle.overflowMode\">ellipsisNone</option><option name=\"charting.axisLabelsX.majorLabelStyle.rotation\">0</option><option name=\"charting.axisTitleX.visibility\">collapsed</option><option name=\"charting.axisTitleY.text\">HTTP status codes</option><option name=\"charting.axisTitleY.visibility\">visible</option><option name=\"charting.axisTitleY2.visibility\">visible</option><option name=\"charting.axisX.abbreviation\">none</option><option name=\"charting.axisX.scale\">linear</option><option name=\"charting.axisY.abbreviation\">none</option><option name=\"charting.axisY.scale\">linear</option><option name=\"charting.axisY2.abbreviation\">none</option><option name=\"charting.axisY2.enabled\">0</option><option name=\"charting.axisY2.scale\">inherit</option><option name=\"charting.chart\">column</option><option name=\"charting.chart.bubbleMaximumSize\">50</option><option name=\"charting.chart.bubbleMinimumSize\">10</option><option name=\"charting.chart.bubbleSizeBy\">area</option><option name=\"charting.chart.nullValueMode\">connect</option><option name=\"charting.chart.showDataLabels\">none</option><option name=\"charting.chart.sliceCollapsingThreshold\">0.01</option><option name=\"charting.chart.stackMode\">default</option><option name=\"charting.chart.style\">shiny</option><option name=\"charting.drilldown\">none</option><option name=\"charting.layout.splitSeries\">0</option><option name=\"charting.layout.splitSeries.allowIndependentYRanges\">0</option><option name=\"charting.legend.labelStyle.overflowMode\">ellipsisMiddle</option><option name=\"charting.legend.mode\">standard</option><option name=\"charting.legend.placement\">right</option><option name=\"charting.lineWidth\">2</option><option name=\"trellis.enabled\">0</option><option name=\"trellis.scales.shared\">1</option><option name=\"trellis.size\">small</option><option name=\"trellis.splitBy\">_aggregation</option></chart></panel></row></dashboard>"

  acl {
	owner = "admin"
	app = "search"
  }
}
```

## Argument Reference
For latest resource argument reference: https://docs.splunk.com/Documentation/Splunk/8.1.1/RESTREF/RESTknowledge#data.2Fui.2Fviews

This resource block supports the following arguments:
* `name` - (Required) Dashboard name.
* `eai:data` - (Required) Dashboard XML definition.

## Attribute Reference
In addition to all arguments above, This resource block exports the following arguments:

* `id` - The ID of the dashboard

## Import

Dashboards/views in the default namespace can be imported by name:

```
terraform import splunk_data_ui_views.example "<view-name>"
```

Dashboards/views in a specific Splunk namespace can be imported with a Splunk REST path or URL. URL-encode the view name when it contains spaces or other special characters:

```
terraform import splunk_data_ui_views.example "/servicesNS/<owner>/<app>/data/ui/views/<url-encoded-view-name>"
```

### After import

Import sets the resource `id` and `name` to the view name. REST-path imports also set initial ACL namespace values (`owner`, `app`, and an inferred `sharing` value). Import does not load every Splunk setting or permission list.

Run `terraform plan` immediately after import. Plan output commonly includes drift until your configuration matches Splunk:

- **ACL drift** — `acl.read` and `acl.write` are not populated during import and may differ from Splunk until you copy values from the Splunk UI or REST API into your `.tf` file. REST-path import infers `sharing` as `app` when `owner` is `nobody`, otherwise `user`. Globally shared views (`sharing = "global"`) may show a one-time ACL change in plan; set `sharing = "global"` explicitly in config if needed.
- **Unset attributes** — Dashboard attributes such as `eai_data` and ACL permissions may differ from a minimal import config until you define them explicitly or use `lifecycle { ignore_changes = [...] }`.
- **Bare-name import** — Importing by name alone does not set `acl`. Without an `acl` block in config, refresh may use provider defaults that do not match the view's Splunk namespace. Prefer REST-path import or set `acl` explicitly.

Recommended workflow:

1. `terraform import ...` (name or REST path)
2. `terraform plan` — review drift
3. Update `.tf` to match required settings, or use `terraform plan -generate-config-out=generated.tf` (Terraform 1.5+) as a starting point
4. Run `terraform plan` again until only intentional changes remain
