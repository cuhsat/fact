# fkey
Shows all [BitLocker Recovery Key IDs](https://learn.microsoft.com/en-us/windows/security/operating-system-security/data-protection/bitlocker/recovery-overview) of an image.

```sh
# fkey [-hv] IMAGE
```

Available options:

- `-h` Show usage
- `-v` Show version

Required system commands:

- [dislocker](https://github.com/Aorimn/dislocker)
- [losetup](https://man7.org/linux/man-pages/man8/losetup.8.html)
- [lsblk](https://man7.org/linux/man-pages/man8/lsblk.8.html)

---
Part of the [Forensic Artifacts Collecting Toolkit](../README.md).