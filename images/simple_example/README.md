# Example Tutorial - AWE Image

This example will show how you can create an AWE-Image, which can be used on the awe learning platform.
We are going to create a simple sql injection with php and mysql.

# Step 1 - Prepare all the needed files

We need two files for our image. The php files containing the sql injection and a init script for our database.
## PHP:
Our php file simple connects to a database and executes a sql query containing a string directly passed into the query (SQL injection!).
```php
<?php
$servername = "127.0.0.1";
$username = "example_user";
$password = "example_password";
$database = "example";

if(isset($_POST["id"])){
    $conn = mysqli_connect($servername, $username, $password,$database);
    
        if(! $conn ) {
            die('Could not connect: ' . mysqli_error());
        }

        $sql = 'SELECT name FROM user where id = '.$_POST["id"];
        $result = mysqli_query($conn, $sql);

        if (mysqli_num_rows($result) > 0) {
            while($row = mysqli_fetch_assoc($result)) {
                echo "Name: " . $row["name"]. "<br>";
            }
        } else {
            echo "0 results";
        }
        mysqli_close($conn);
    }
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Example</title>
</head>
<body>
    <h1>Search for User by ID:</h1>
    <form action="/" method="POST">
        <input type="text" id="name" name="id">
        <button type="submit">Search ...</button>
    </form>
</body>
</html>
```

## init.sql:
Our init.sql creates a new user and database. It sets the permissions for the user on the database and creates the needed tables.
It also inserts some example data into the tables.
(And gives the users the file privilege ... to make the sql injection a litte bit more interessting)
```sql
CREATE USER 'example_user'@'localhost' IDENTIFIED BY 'example_password';
CREATE DATABASE example;
GRANT ALL PRIVILEGES ON example.* TO 'example_user'@'localhost';
GRANT FILE ON *.* TO 'example_user'@'localhost';
USE example;

CREATE TABLE user (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name varchar(255)
);

INSERT INTO user(name) VALUES ('john'),('harry'),('rite'),('Sophie');

CREATE TABLE secret (
    secret varchar(255)
);

INSERT INTO secret(secret) VALUES ('th!5!sSup3rS3cr3T!');
```

# Step 2

We know that we are going to use php and mysql. We now have two options. We could create our image using a Linux distribution base image like Ubuntu, Alpine or Debian, but we can also use a image which contains a little more software.
In this example, we are using the PHP base image. It is based on Debian, but comes preinstalled and preconfigured with PHP. This makes our Dockerfile a little bit easier.

```Dockerfile
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
```

# Step 3 - Building the image
When we have everything prepared, we only need to build our image. We can archive this by running the following command.

```bash
docker build -t example01 .
```

The `-t example01` argument is necessary and gives the image a name/tag.
After some time the image is build. We can see it in our docker images with the following command.

```bash
$ docker image ls
REPOSITORY          TAG       IMAGE ID       CREATED        SIZE
example01           latest    7ff103d8f634   2 hours ago    846MB <--- HERE IT IS!
hackme02            latest    44a3264eaa8c   5 weeks ago    534MB
hackme01            latest    7dcb2b064742   5 weeks ago    942MB
docker101tutorial   latest    80b65cb2639a   8 weeks ago    27.7MB
alpine/git          latest    ed0ba0fc6585   2 months ago   28.4MB
alpine              latest    389fef711851   2 months ago   5.58MB
```

# Step 4 - Export the image

When we are building the image on the same machine where awe is running, the image will become available to the learing platform automatically.
If we are building the image on another machine, we need to export it, and import it into the awe platform.

To export the image we run the following command:
```bash
docker save example01 > image.tar
```
`example01` is the tag (or name) of our image.

# Step 5 - Import the image

To import the image into an awe instance, open the admin page and upload the created tar file from the previous step using the admin password. Done.

# Step 6 - Exploit 😈

When we start the image through our learning platform and enter the site we are greeted with a form. We know that this form contains a sql injection. Here are three examples how we can exploit this:

```sql
-- bypass the filter
1 or 1=1

-- reading from another table
1 union select secret from secret

-- reading a file .. also works with /etc/passwd
1 union all select load_file("/flag.txt")
```