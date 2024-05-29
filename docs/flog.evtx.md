# flog.evtx
Log [Windows event logs](https://forensics.wiki/windows_event_log_%28evt%29/) as JSON in [ECS](https://www.elastic.co/guide/en/ecs/current/index.html).

```sh
$ flog.evtx [-pqhv] [-D DIRECTORY] [FILE ...]
```

Available options:

- `-D` Log directory
- `-p` Pretty JSON
- `-q` Quiet mode
- `-h` Show usage
- `-v` Show version

Required system commands:

- [dotnet](https://dotnet.microsoft.com/en-us/download/dotnet/6.0)

> Use `make tools` to install [Eric Zimmerman's Tools](https://ericzimmerman.github.io/#!index.md).

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).