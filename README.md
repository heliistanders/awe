# AWE - Advanced Web Exploitation - WIP!
AWE is a learning platform for advanced web exploitation technics. This project will include vulnerabilties on purpose - please don't run this on a machine which can be accessed by others - run this in a VM! 

## Architecture
AWE is designed to run in a Linux VM. AWE runs a go application which controls the docker deamon. The single exploit targets are provided as docker images. AWE manages the single images.

# Requirements
- docker
- a user which can access docker  
- golang (for development)

# Installation & 

```bash
$ git clone https://github.com/heliistanders/awe
$ cd awe
$ go build
$ ./awe
```

## Creating an AWE Docker Image:

Creating a docker image for AWE is the same as any other docker image. The only difference is that are requires some labes to be set.
- awe=NAME
- difficulty=DIFFICULTY
- ports=PORT1[,PORT2]

The awe label provides the name of the machine and gets displayed to the user.

The difficulty label provides the expected difficult of the machine and gets displayed to the user.

The ports label contains one port or a comma-separated list of ports. Those ports are used internally by the services running inside the container and get mapped to random ports accessible to the user.

### Example Build command when using a Dockerfile

Add the labels when building the image

```bash
 docker build -t hackme01 --label awe="Hackme 01" --label difficulty=easy --label ports=80 --label hint="Deserialization" .
 docker build -t hackme02 --label awe="Hackme 02" --label difficulty=easy --label ports=80,8080 .
 ```

OR provide the labels inside the Dockerfile

```dockerfile
FROM php:8.0-apache
LABEL awe="Hackme 01"
LABEL ports="80"
LABEL difficulty="easy"
LABEL hint="Deserialization"
COPY ./www/hackme/ /var/www/hackme
COPY ./www/html/ /var/www/html/
COPY ./opt/vuln/ /opt/vuln/
COPY ./hackme.fhj.conf /etc/apache2/sites-available/
RUN apt update; \
	mkdir -p /usr/share/man/man1; \
	apt install -y jruby; \
	cd /etc/apache2/sites-available & a2ensite hackme.fhj.conf
```
when provided inside the Dockerfile, build it with:

```bash
docker build -t hackme01 .
```
## License

MIT - see [LICENSE](./LICENSE) for further information

## ToDo

- [x] Startup - Check docker, database ...
- [x] Serve static content
- [x] Rework machine handling (restart, helper functions)
- [ ] Implement better logging  
- [ ] Web terminal via WebSocket into solved machines
- [x] Refactor Codebase (especially the database handling)
- [ ] Add Frontend as git submodule
- [x] Upload AWE-Docker Images
- [ ] Prevent everyone from uploading an Image (otherwise the pc ca be taken over)
- [x] Additional flags for hints?