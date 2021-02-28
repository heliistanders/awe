FROM php:8.0-apache
LABEL awe="Hackme 01"
LABEL difficulty="easy"
LABEL hint="Deserialization"
EXPOSE 80
COPY ./www/hackme/ /var/www/hackme
COPY ./www/html/ /var/www/html/
COPY ./opt/vuln/ /opt/vuln/
COPY ./hackme.fhj.conf /etc/apache2/sites-available/
RUN apt update; \
	mkdir -p /usr/share/man/man1; \
	apt install -y jruby; \
	cd /etc/apache2/sites-available & a2ensite hackme.fhj.conf
