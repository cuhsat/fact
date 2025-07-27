# Forensic Artifacts Collecting Toolkit

A basic shell pipeline for extracting forensic artifacts from disk images. Relevant artifacts will be processed and provided in [ECS](https://www.elastic.co/guide/en/ecs/current/index.html) format for ingestion with [Logstash](https://www.elastic.co/de/logstash).

```console
# fmount disk.dd | ffind | flog -D artifacts
```

## Tools
- [fmount](docs/fmount.md)
- [ffind](docs/ffind.md)
- [flog](docs/flog.md)

## License
Released under the [MIT License](LICENSE.md).