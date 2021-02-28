# AWE - Advanced Web Exploitation - WIP!
AWE is a learning platform for advanced web exploitation technics. 
This project will is indented to serve vulnerabilities.
Please don't run this on a machine which can be accessed by others - run this in a VM! 

## Architecture
AWE is designed to run in a Linux VM. Executables for macOS and windows can be compiled, but are not tested.
AWE orchestrated the docker daemon and provided an API to manage the learning platform.
AWE itself is written in GO. It uses SQLite3 to persist data.
A frontend written with vue can be found here: https://github.com/heliistanders/awe-frontend

# Requirements
- docker
- a user which can access docker (root, or a user in the docker group)
- golang (for development)

# Installation

```bash
$ git clone https://github.com/heliistanders/awe
$ cd awe
$ go build
$ ./awe
or
$ AWE_PASS=<admin password> ./awe
```

To upload a new challenge a password is needed. The password can be set via environment variable.
```bash
$ AWE_PASS=better_secret ./awe
```
If no environment variable is provided, the password defaults to a random string.
This string gets printed when starting the app.

## Creating an AWE Docker Image:

Creating a challenge for the awe platform is the same as creating a docker image. The only difference is, that
an awe challenge needs 2-3 labels to be set.
- awe=NAME
- difficulty=DIFFICULTY
- hint=<hint> (optional)

The awe label provides the name of the challenge and gets displayed to the user.
The difficulty label provides the expected difficult of the challenge and gets displayed to the user.

### PORTS
It is important to set the exposed ports. The awe platform needs to know which ports to forward.

### Example Challenges

You can find example challenges in the examples folder. It shows how to build then. Start with the simple_example.

If you are building the challenges on the same host, as awe is running, the challenges are automatically picked up by the platform.
If you are bulding the challenges on another host, you can export them via tar file and add them via the admin page of the front-end.


## License

MIT - see [LICENSE](./LICENSE) for further information

## ToDo

- [x] Startup - Check docker, database ...
- [x] Serve static content
- [x] Rework machine handling (restart, helper functions)
- [x] Implement better logging  
- [x] Web terminal via WebSocket into solved machines
- [x] Refactor Codebase (especially the database handling)
- [x] Upload AWE-Docker Images
- [x] Prevent everyone from uploading an Image (otherwise the pc ca be taken over)
- [x] Additional flags for hints?