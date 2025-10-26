# conservetp
dockerized microservice conservetp


## Docker

### Compose
`develop.watch{action: rebuild, path: ./, target: /app}`

Command to rebuild the container after changes to the path are made. Target path is the location within the docker container rebuild is to take place.

- Use rebuild to recreate the docker image and re-initialize the go server. 
- sync+restart is available but doesn't build the container again (issue with go needing to be compiled.)
    - Must use a two stage rebuild (build binary then cmd ["./main"])