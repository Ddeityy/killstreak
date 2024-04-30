# Better killstreaks for TF2

## Provides a more informative format for killstreaks.

### TF2 demo [parser](https://github.com/demostf/parser) by [@icewind1991](https://github.com/icewind1991/).

### Old _events.txt
```
>
[2023/11/08 23:48] Bookmark ("2023-11-08_23-32-45" at 20000)
[2023/11/08 23:48] Killstreak 5 ("2023-11-08_23-32-45" at 60781)
[2023/11/08 23:48] Killstreak 6 ("2023-11-08_23-32-45" at 62998)
>
```

### New events.txt (playdemo starts 500 ticks before killstreak/bookmark)
```
>
[2023-11-05] pl_upward_f10 sniper
[2023-11-05] Killstreak 5 38690-40589 [28.48 seconds]            playdemo demos/2023-11-05_23-10-39; demo_gototick 38190 0 1
[2023-11-05] Bookmark at 41189                                   playdemo demos/2023-11-05_23-10-39; demo_gototick 40689 0 1
>
```
### Install
* Download the linux release and run
```console
$ sudo ./install.sh
```

### processDemos
Takes your existing _events.txt and parses all existing demos for killstreaks + bookmarks (if your ds_mark comment contains "bookmark") to events.txt
```console
$ ./processDemos
```
