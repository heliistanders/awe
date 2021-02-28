###########################
#                         #
#    AWE Example Image    #
#                         #
###########################

# This line is a comment. All comments can be removed without changing the function of the dockerfile.

# We are using the PHP Base Image
FROM php:8-apache
# Our Image will expose port 80, the port apache is running on. This is a necessary information for AWE!
EXPOSE 80
# Run commands are executed at build time. We are installing all the things we need.
# Update our apt repository
RUN apt update
# Install mysql. Dont forget to use the '-y' switch, to automatically accept the installation
RUN apt -y install default-mysql-server
# Install the mysqli extension. This allow us to connect to our mysql server from php
RUN docker-php-ext-install mysqli
# Copy our files into the image COPY <src> <dst>
COPY ./index.php /var/www/html
COPY ./init.sql /init.sql

# Those are AWE specific labels
# The awe=<name> label MUST be set
LABEL awe="Example"
# the other two are optional, but provide information to our users
LABEL difficulty="easy"
LABEL hint="SQL Injection"

# This gets run when we are starting a container. In this example, we are starting
# mysql and apache. Then we set the database using our init.sql script. After that we
# are running a command which does not end, so the container keeps running
CMD service mysql start && service apache2 start && mysql < /init.sql && tail -F /var/log/mysql/error.log