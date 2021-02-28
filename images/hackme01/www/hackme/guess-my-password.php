<?php

if(isset($_POST['data'])){
  $guess = hash('sha256', $_POST['data'], false);
  if(strcmp($guess, 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3') == 0){
    $output = "YOU ARE RIGHT!";
  } else {
    $output = "YOU ARE WRONG! OR A BAD HACKER!";
  }
}
?>
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Online JSON parser</title>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">

	<style>
		html, body {
			height: 100%;
			width: 100%;
			padding: 0;
			margin: 0;
			background-color: rgba(0,0,0,0.4);
			display: flex;
			justify-content: center;
			align-items: center;
		}

		input {
			box-shadow: 0px 1px 2px rgba(0,0,0,.4);
			border-radius: 5px;
			border: 1px solid black;
			width: 100%;
			height: 48px;
			padding: 0px 12px;
		}

		h1 {
			text-align: center;
		}

		form {
			display: flex;
			flex-direction: column;
			align-items: center;
		}
	</style>
	</head>
	<body>
		<form  action="" method="post">
			<h1 >
				Guess my password!
			</h1>
				<input  type="text" name="data"></input>
			<pre>
				<?php if(isset($_POST['data'])) { echo $output; }?>
			</pre>
			<button  type="submit">
				GUESS
			</button>
		</form>
		
	</body>
</html>
