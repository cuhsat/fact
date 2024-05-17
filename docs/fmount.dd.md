# fmount.dd
Mount forensic raw or dd disk images for read-only processing.

```sh
$ fmount.dd [-fsuzhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-D DIRECTORY] IMAGE
```

Available options:

- `-D` Mount point
- `-H` Hash algorithm
- `-V` Verify hash sum
- `-f` Force type
- `-s` System partition only
- `-u` Unmount image
- `-z` Unzip image
- `-h` Show usage
- `-v` Show version

Required system commands:

- [losetup](https://man7.org/linux/man-pages/man8/losetup.8.html)
- [lsblk](https://man7.org/linux/man-pages/man8/lsblk.8.html)
- [mount](https://man7.org/linux/man-pages/man8/mount.8.html)
- [umount](https://man7.org/linux/man-pages/man8/umount.8.html)

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).