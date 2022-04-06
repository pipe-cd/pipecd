# piped-base-okd
This directory contains the Piped docker image packing with [nss_wrapper](https://cwrap.org/nss_wrapper.html) to add the logged-in UID to "passwd" at runtime without having to directly modify `/etc/passwd`.
This mainly aims to workaround to deal with some issues due to random UID on OpenShift less than 4.2.

For more why it got needed, see: https://github.com/pipe-cd/pipecd/issues/1905
