<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Moo-Say</title>
  <style>
    html,body {
      height: 100%;
      width: 100%;
    }

    body {
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .main {
      display: block;
      width: 500px;
    }

    form {
      text-align: center;
    }

    pre {
      font-size: 24px;
    }
  </style>
</head>
<body>
  <div class="main">

    <form action="moo-say.php" method="post">
      <input type="text" name="what" id="what">
      <button type="submit">SAY!</button>
    </form>
    <pre>
      <?php
    $text = "srsly dude, why?";
    if(isset($_POST["what"])){
      if(!empty($_POST["what"])){
        if($_POST["what"] == "'"){
          $text = "1064 - You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ')'' at line 1 ..... just kidding there's no sqli :)";
        } else {

          $text = $_POST["what"];
          if(strlen($text) > 16){
            $text = chunk_split($_POST["what"], 16, " &gt;\n    &lt; ");
            
          } else {
            $text = str_pad($text, 16)." >";
          }
          $text = htmlspecialchars($text);
        }
      }
    }
    echo "------------------
     < $text >
      ------------------
            \   ^__^
             \  (oo)\_______
                (__)\       )\/\
                    ||----w |
                    ||     ||";
    
    ?>

</pre>
</div>
</body>
</html>
