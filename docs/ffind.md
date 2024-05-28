# ffind
Find forensic artifacts in mount points or on the live system.

```sh
$ ffind [-rsuqhv] [-H CRC32|MD5|SHA1|SHA256] [-Z ARCHIVE] [-L FILE] [MOUNT ...]
```

Available options:

- `-H` Hash algorithm
- `-Z` Archive name
- `-L` Listing name
- `-r` Relative paths
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

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).