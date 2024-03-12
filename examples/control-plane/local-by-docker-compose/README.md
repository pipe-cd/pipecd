# docker-compose

This example shows how to install a Control Plane on local by docker-compose.

1. Install
```sh
docker-compose up
```

2. Clean up
```sh
docker-compose down
```

NOTE: By following commands instead of above `down`, you can keep data such as Piped or applications on the Control Plane even after restarting/updating the server component.

```sh
# Restart only the server component.
docker-compose rm -fsv pipecd-server
docker-compose up pipecd-server
```
