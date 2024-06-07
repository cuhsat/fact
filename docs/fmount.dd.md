# fmount.dd
Mount forensic raw or dd disk images for read-only processing.

```sh
# fmount.dd [-fruszqhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-B KEY] [-D DIR] IMAGE
```

Available options:

- `-D` Mount point
- `-B` BitLocker key
- `-H` Hash algorithm
- `-V` Verify hash sum
- `-f` Force type
- `-r` Recovery key ids 
- `-u` Unmount image
- `-s` System partition only
- `-z` Unzip image
- `-q` Quiet mode
- `-h` Show usage
- `-v` Show version

Required system commands:

- [dislocker](https://github.com/Aorimn/dislocker)
- [losetup](https://man7.org/linux/man-pages/man8/losetup.8.html)
- [lsblk](https://man7.org/linux/man-pages/man8/lsblk.8.html)
- [mount](https://man7.org/linux/man-pages/man8/mount.8.html)
- [umount](https://man7.org/linux/man-pages/man8/umount.8.html)

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).