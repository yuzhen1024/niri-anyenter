sudo chgrp input $GOBIN/type-allthing-run
sudo chmod g+s $GOBIN/type-allthing-run
sudo setcap cap_setgid+ep $GOBIN/type-allthing-run
