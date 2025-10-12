## install

`go install github.com/yuzhen1024/niri-anyenter`


## build install

```bash
pacman -S systemd libinput pkgconf zig

git clone https://github.com/yuzhen1024/niri-anyenter.git
cd niri-anyenter/c
zig build
sh _move-output.sh
cd ..
go install
sh _install_chmod.sh
```
