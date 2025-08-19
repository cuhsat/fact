# Forensic Artifacts Collecting Toolkit

A shell pipeline for extracting forensic artifacts from disk images in [ECS](https://www.elastic.co/guide/en/ecs/current/index.html) format. Important artifacts will be processed and provided for ingestion with [Logstash](https://www.elastic.co/de/logstash).

```console
# fmount disk.raw | ffind | flog -D logstash
```

## [fmount](https://github.com/cuhsat/fmount)
Mount various disk images for forensic read-only processing.

## [ffind](https://github.com/cuhsat/ffind)
Find forensic artifacts in mount points or the live system.

## [flog](https://github.com/cuhsat/flog)
Log forensic artifacts as JSON in ECS format.

## License
All released under the [MIT License](LICENSE.md).
