sudo chgrp input $GOBIN/niri-anyenter
sudo chmod g+s $GOBIN/niri-anyenter
sudo setcap cap_setgid+ep $GOBIN/niri-anyenter
