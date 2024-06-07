# fmount
Mount forensic disk images for read-only processing.

```sh
# fmount [-suzqhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-B KEY] [-T RAW|DD|VMDK] [-D DIR] IMAGE
```

Available options:

- `-D` Mount point
- `-T` Image type
- `-B` BitLocker key
- `-H` Hash algorithm
- `-V` Verify hash sum
- `-s` System partition only
- `-u` Unmount image
- `-z` Unzip image
- `-q` Quiet mode
- `-h` Show usage
- `-v` Show version

Supported disk formats:

- [VMDK](https://forensics.wiki/vmware_virtual_disk_format_%28vmdk%29/)
- [DD (Raw)](https://forensics.wiki/raw_image_format/)

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).