# Dell OS10 Exporter

Prometheus exporter for Dell OS10 switches. Uses the 
[SmartFabric OS10 RESTCONF API ](https://developer.dell.com/apis/13448/versions/10.5.4.0/docs/introduction.md) 
to query VLT status,port-channel status, and power supply status. This exporter
is intended to query multiple Dell OS10 switches from an external host.

The `/dellos10` metrics endpoint exposes the Dell OS10 metrics and requires a
`target` parameter.  The `module` parameter can also be used to select which
probe commands to run, the default module is `power`. Available modules are:

- system
- vltlocalinfo
- vltpeerinfo
- vltportchannel

The `/metrics` endpoint exposes Go and process metrics for this exporter.

## Configuration

This exporter requires a dellos10.conf file. Example config:

```ini
[connection:switch1.example.com]
host=switch1.example.com
username=admin
password=root
enablepwd=passwd

[connection:switch2.example.com]
host=switch2.example.com
username=admin
password=root
enablepwd=passwd
```

## Prometheus configs

```yaml
- job_name: dellos10
  metrics_path: /dellos10
  static_configs:
  - targets:
    - switch1.example.com
    - switch2.example.com
    labels:
      module: system,vltportchannel,vltpeerinfo,vltlocalinfo
  - targets:
    - switch3.example.com
    labels:
      module: system,vltportchannel
  relabel_configs:
  - source_labels: [__address__]
    target_label: __param_target
  - source_labels: [__param_target]
    target_label: instance
  - source_labels: [module]
    target_label: __param_module
  - target_label: __address__
    replacement: 127.0.0.1:9465
```

Example systemd unit file [here](systemd/dellos10_exporter.service)

## Sample Metrics

```
dellos10_power_supply_state{powerSupply="1",state="up"} 1
dellos10_power_supply_state{powerSupply="2",state="up"} 1
dellos10_vlt_unit_id{role="primary"} 2
dellos10_vlt_peer_info{role="secondary",status="up",unitId="1"} 1
dellos10_vlt_portchannel_ports{interface="port-channel1",state="active"} 0
dellos10_vlt_portchannel_ports{interface="port-channel1",state="inactive"} 2
dellos10_vlt_portchannel_ports{interface="port-channel2",state="active"} 1
dellos10_vlt_portchannel_ports{interface="port-channel2",state="inactive"} 1
dellos10_vlt_portchannel_ports{interface="port-channel2",state="active"} 2
dellos10_vlt_portchannel_ports{interface="port-channel2",state="inactive"} 0
```

## License

dellos10_exporter is released under the Apache License Version 2.0. See the LICENSE file.
