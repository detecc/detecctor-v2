# Plugin service configuration

## Configuration file

Before the Detecctor Plugin Service is run, it _must_ have a configuration file set up. The attributes, required to
successfully run the service, are:

1. mqttBroker information
2. plugin directory and list
3. database information

The configuration file formats supported are **YAML, JSON and TOML**. An example `config` file in the`yaml` format:

```yaml
mqttBroker:
  host: localhost
  port: 7777
  username: "username"
  password: "pass"
  clientId: "clientId1"
  tls:
plugins:
  dir: ../detecc-core/compiled/server
  list:
    - "hw-monitor"
mongo:
  database: "test"
  host: localhost
  username: root
  password: pass12345
  port: 27017

```

## Flags

To change the location of the configuration files, to enable persistence, both the configuration file location and
format can be specified with flags.

```bash
./main --help 
    --config-format # Format of the configuration files (yaml, json or toml)
    --config-file # Path of the configuration file (default: working directory)
```
