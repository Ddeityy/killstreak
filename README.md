# Better killstreaks for TF2

## Provides a more informative format for killstreaks.

### TF2 demo [parser](https://github.com/demostf/parser) and [cutter](https://github.com/demostf/edit) by [@icewind1991](https://github.com/icewind1991/).

### Old
```
[2023/11/08 23:48] Bookmark ("2023-11-08_23-32-45" at 20000)
[2023/11/08 23:48] Killstreak 5 ("2023-11-08_23-32-45" at 60781)
[2023/11/08 23:48] Killstreak 6 ("2023-11-08_23-32-45" at 62998)
```

### New
```
cut=true
playdemo demos/cut_2023-11-08_23-32-45; demo_gototick 60781 0 1
cut=false
playdemo demos/2023-11-08_23-32-45; demo_gototick 60781 0 1

[2023/11/08 23:48] cp_entropy_b5 scout
[2023/11/08 23:48] Bookmark ("2023-11-08_23-32-45" at 20000)
[2023/11/08 23:48] Killstreak 6 ("2023-11-08_23-32-45" 60781-62998 [18.25 seconds])
```
### Install
#### Linux:
* Download the linux release and run
```console
$ sudo ./install.sh
```

#### Windows
* Download the windows release
* Download [nssm](http://nssm.cc/download)
* Run
```console
nssm.exe install killstreak
```
Select the main killstreak.exe path and pass cut=true in the arguements for automatic demo cutting