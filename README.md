# Better killstreaks for TF2 for Linux

## Provides a more informative format for killstreaks when using TF2's default demo recorder.

### TF2 demo [parser](https://github.com/demostf/parser) by [@icewind1991](https://github.com/icewind1991/).

### Old
```
[2023/11/08 23:48] Bookmark ("2023-11-08_23-32-45" at 20000)
[2023/11/08 23:48] Killstreak 5 ("2023-11-08_23-32-45" at 60781)
[2023/11/08 23:48] Killstreak 6 ("2023-11-08_23-32-45" at 62998)
```

### New
```
[2023/11/08 23:48] cp_entropy_b5 scout
[2023/11/08 23:48] Bookmark ("2023-11-08_23-32-45" at 20000)
[2023/11/08 23:48] Killstreak 6 ("2023-11-08_23-32-45" 60781-62998 [18.25 seconds])
```
#### **For now this only works with the default demo folder location: tf/demos**
### Installation
Download the release and run
```console
$ sudo ./install.sh
```
