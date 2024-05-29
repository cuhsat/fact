# flog
Log forensic artifacts as JSON in [ECS](https://www.elastic.co/guide/en/ecs/current/index.html).

```sh
$ flog [-hv] [-D DIRECTORY] [FILE ...]
```

Available options:

- `-D` Log directory
- `-h` Show usage
- `-v` Show version

Supported artifacts for Windows 7+ systems:

- [System Event Logs](flog.evtx.md)

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).