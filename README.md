## install

```bash
go install github.com/yuzhen1024/niri-anyenter
sudo chgrp input $GOBIN/niri-anyenter
sudo chmod g+s $GOBIN/niri-anyenter
sudo setcap cap_setgid+ep $GOBIN/niri-anyenter
```

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

## usage

```kdl
spawn-at-startup "PATH_YOUR_GO_BIN/niri-anyenter"
```

```txt
usage: niri-anyenter [<flags>]

Flags:
  --[no-]help                 Show context-sensitive help (also try --help-long and --help-man).
  --launcher="fuzzel"         launcher or runner, run your apps, e.g. fuzzel
  --mode=preinput             default mode is preinput mode, in future updates include uinput mode(maybe), uinput mode is wip,
                              now only use preinput.
  --preinput-flag="--search"  in empty workspace, type allthing start launcher, now need set for preinput whatever in this.
                              for fuzzel is --search, for rofi is -filter, you can `--help` your using launcher for check.
  --preinput-delay-ms=400     input anything run launcher, but how quick can you type? need the this value to check if you
                              stopped typing.
  --exclude-windows= ...      gnore windows if existing, the windows will not stop run launcher. example: --exclude-windows
                              "app1" "app2" ...
  --exclude-layers=nwg-drawer ...  
                              this will stop launcher run, if layers existing. example: --exclude-layers "app1" "app2" ...
  --launcher-rule="layer=launcher"  
                              use `sleep 3 && niri msg layers`see namespace to check. the fuzzel use layer, the rofi use
                              window, so, set you want. default value is `layer=launcher`, fit for fuzzel.
  --[no-]version              Show application version.
```
