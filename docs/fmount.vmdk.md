# fmount.vmdk
Mount forensic [VMDK](https://forensics.wiki/vmware_virtual_disk_format_%28vmdk%29/) disk images for read-only processing.

```sh
# fmount.vmdk [-fruszqhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-B KEY] [-D DIR] IMAGE
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
- [qemu-nbd](https://www.qemu.org/docs/master/tools/qemu-nbd.html)
- [lsblk](https://man7.org/linux/man-pages/man8/lsblk.8.html)
- [mount](https://man7.org/linux/man-pages/man8/mount.8.html)
- [umount](https://man7.org/linux/man-pages/man8/umount.8.html)

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).