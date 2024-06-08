# FACT
Forensic Artifacts Collecting Toolkit.

```sh
# fmount image.dd | ffind | flog
```

## Tools

### fmount
Mount disk images for read-only processing.

```sh
# fmount [-ruszqhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-B KEY] [-D DIR] IMAGE
```

Available options:

- `-D` Mount point
- `-B` BitLocker key
- `-H` Hash algorithm
- `-V` Verify hash sum
- `-r` Recovery key ids 
- `-u` Unmount image
- `-s` System partition only
- `-z` Unzip image
- `-q` Quiet mode
- `-h` Show usage
- `-v` Show version

Supported image types on Linux systems:

- [vdi](https://forensics.wiki/virtual_disk_image_%28vdi%29/)
- [vpc](https://cloud.ibm.com/docs/vpc?topic=vpc-planning-custom-images)
- [vhdx](https://forensics.wiki/virtual_hard_disk_%28vhd%29/)
- [vmdk](https://forensics.wiki/vmware_virtual_disk_format_%28vmdk%29/)
- [parallels](https://github.com/libyal/libphdi/blob/main/documentation/Parallels%20Hard%20Disk%20image%20format.asciidoc)
- [qcow2](https://forensics.wiki/qcow_image_format/)
- [qcow](https://forensics.wiki/qcow_image_format/)
- [raw](https://forensics.wiki/raw_image_format/)

Required system commands:

- [dislocker](https://github.com/Aorimn/dislocker)
- [qemu-nbd](https://www.qemu.org/docs/master/tools/qemu-nbd.html)
- [lsblk](https://man7.org/linux/man-pages/man8/lsblk.8.html)
- [mount](https://man7.org/linux/man-pages/man8/mount.8.html)
- [umount](https://man7.org/linux/man-pages/man8/umount.8.html)

### ffind
Find forensic artifacts in mount points or on the live system.

```sh
$ ffind [-rcsuqhv] [-H CRC32|MD5|SHA1|SHA256] [-C CSV] [-Z ZIP] [MOUNT ...]
```

Available options:

- `-H` Hash algorithm
- `-C` CSV listing name
- `-Z` Zip archive name
- `-r` Relative paths
- `-c` Volume shadow copy
- `-s` System artifacts only
- `-u` User artifacts only
- `-q` Quiet mode
- `-h` Show usage
- `-v` Show version

Supported artifacts for Windows 7+ systems:

- [System Active Directory](https://forensics.wiki/active_directory/)
- [System Registry Hives](https://forensics.wiki/windows_registry/)
- [System Prefetch Files](https://forensics.wiki/prefetch/)
- [System Event Logs](https://forensics.wiki/windows_event_log_%28evt%29/)
- [System AmCache](https://forensics.wiki/amcache/)
- [User Registry Hives](https://forensics.wiki/windows_registry/)
- [User Jump Lists](https://forensics.wiki/jump_lists/)
- [User Browser Histories](https://forensics.wiki/google_chrome/)

### flog
Log forensic artifacts as JSON in [ECS](https://www.elastic.co/guide/en/ecs/current/index.html).

```sh
$ flog [-pqhv] [-D DIRECTORY] [FILE ...]
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

#### Roadmap
- [ ] Support for [System Active Directory](https://forensics.wiki/active_directory/)
- [ ] Support for [System Registry Hives](https://forensics.wiki/windows_registry/)
- [ ] Support for [System Prefetch Files](https://forensics.wiki/prefetch/)
- [x] Support for [System Event Logs](https://forensics.wiki/windows_event_log_%28evt%29/)
- [ ] Support for [System AmCache](https://forensics.wiki/amcache/)
- [ ] Support for [User Registry Hives](https://forensics.wiki/windows_registry/)
- [ ] Support for [User Jump Lists](https://forensics.wiki/jump_lists/)
- [ ] Support for [User Browser Histories](https://forensics.wiki/google_chrome/)

## License
Released under the [MIT License](LICENSE).