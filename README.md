# atlas2sentinel

A Go program that exports MongoDB Atlas logs to Microsoft Sentinel SIEM.

## Running

First create a yaml file, such as `config.yml`:
```yaml
log:
  level: INFO

microsoft:
  app_id: ""
  secret_key: ""
  tenant_id: ""
  subscription_id: ""
  resource_group: ""
  workspace_name: ""

  dcr:
    endpoint: ""
    rule_id: ""
    stream_name: ""

  expires_months: 6
  update_table: false

atlas:
  url: ""
  lookback_days: 1
  public_key: ""
  private_key: ""
```

And now run the program from source code:
```shell
% make
go run ./cmd/... -config=dev.yml
INFO[0000] shipping logs                                 module=sentinel_logs table_name=AtlasLogs total=82
INFO[0002] shipped logs                                  module=sentinel_logs table_name=AtlasLogs
INFO[0002] successfully sent logs to sentinel            total=82
```

Or binary:
```shell
% one2sen -config=config.yml
```

## Building

```shell
% make build
```
