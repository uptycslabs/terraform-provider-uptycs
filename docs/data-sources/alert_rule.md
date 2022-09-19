---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "uptycs_alert_rule Data Source - terraform-provider-uptycs"
subcategory: ""
description: |-
  
---

# uptycs_alert_rule (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `alert_tags` (List of String)
- `code` (String)
- `custom` (Boolean)
- `description` (String)
- `destinations` (Attributes List) (see [below for nested schema](#nestedatt--destinations))
- `enabled` (Boolean)
- `grouping` (String)
- `grouping_l2` (String)
- `grouping_l3` (String)
- `is_internal` (Boolean)
- `lock` (Boolean)
- `name` (String)
- `notify_count` (Number)
- `notify_interval` (Number)
- `rule` (String)
- `rule_exceptions` (List of String)
- `sql_config` (Attributes) (see [below for nested schema](#nestedatt--sql_config))
- `throttled` (Boolean)
- `type` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--destinations"></a>
### Nested Schema for `destinations`

Optional:

- `close_after_delivery` (Boolean)
- `destination_id` (String)
- `notify_every_alert` (Boolean)
- `severity` (String)


<a id="nestedatt--sql_config"></a>
### Nested Schema for `sql_config`

Optional:

- `interval_seconds` (Number)

