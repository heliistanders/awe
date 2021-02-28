# Example02 - Hackme01

This example exposes a serialization vulnerability in the java jackson library.
It uses JRUBY to mock a JAVA environment and uses PHP as a frontend.

```dockerfile
FROM php:8.0-apache
EXPOSE 80
LABEL awe="Hackme 01"
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