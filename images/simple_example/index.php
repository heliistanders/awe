<?php
$servername = "127.0.0.1";
$username = "example_user";
$password = "example_password";
$database = "example";
// Create connection

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