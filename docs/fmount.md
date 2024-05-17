# fmount
Mount forensic disk images for read-only processing.

```sh
$ fmount [-suzhv] [-T RAW|DD] [-D DIRECTORY] IMAGE
```

Available options:

- `-D` Mount point
- `-T` Image type
- `-s` System partition only
- `-u` Unmount image
- `-z` Unzip image
- `-h` Show usage
- `-v` Show version

Supported disk formats:

- [DD (Raw)](https://forensics.wiki/raw_image_format/)

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).