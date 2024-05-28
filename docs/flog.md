# flog
Log forensic artifacts in [ECS](https://www.elastic.co/guide/en/ecs/current/index.html) schema.

```sh
$ flog [-hv] [-D DIRECTORY] [FILE ...]
```

Available options:

- `-D` Log directory
- `-h` Show usage
- `-v` Show version

Supported artifacts for Windows 7+ systems:

- [System Event Logs](flog.evt.md)

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).