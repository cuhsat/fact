# flog.evt
Log Windows event log artifacts in [ECS](https://www.elastic.co/guide/en/ecs/current/index.html) schema.

```sh
$ flog.evt [-hv] [-D DIRECTORY] [FILE ...]
```

Available options:

- `-D` Log directory
- `-h` Show usage
- `-v` Show version

Required system commands:

- [dotnet](https://dotnet.microsoft.com/en-us/download/dotnet/6.0)

> Use `scripts/eztools.sh` to install [Eric Zimmerman's Tools](https://ericzimmerman.github.io/#!index.md).

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).