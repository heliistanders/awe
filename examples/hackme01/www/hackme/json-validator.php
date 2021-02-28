<?php
if(isset($_POST['data'])){
	$error = "";
	$filename = tempnam("/dev/shm", "payload");
	$myfile = fopen($filename, "w") or die("Unable to open file!");
	$txt = $_POST['data'];
	fwrite($myfile, $txt);
	fclose($myfile);
	exec("/usr/bin/jruby /opt/vuln/parse.rb $filename 2>&1", $cmdout, $ret);
	unlink($filename);

	if($ret === 0){
		$output = 'Validation successful!';
	} else {
		$output = 'Validation failed: '. $cmdout[1];
		$error = implode($cmdout);
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

textarea {
	box-shadow: 0px 1px 2px rgba(0,0,0,.4);
	border-radius: 5px;
	border: 1px solid black;
	width: 100%;
	height: 200px;
	padding: 12px 12px;
}

h1 {
	text-align: center;
}

form {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.output {
	margin-top: 20px;
	margin-bottom: 30px;
	padding: 5px 10px;
	max-width: 500px;
	background: white;

	box-shadow: 0px 1px 2px rgba(0,0,0,.4);
	border-radius: 5px;
	border: 1px solid black;
</style>

</head>
<body>
		<form  action="" method="post">
			<h1>
				Online JSON validator! (WIP)
			</h1>
			<textarea  type="text" name="data" cols="50"></textarea>
			<div class="output">
				<?php if(isset($_POST['data'])) { echo $output; } else { echo 'Output goes here!';}?>
			</div>
			<button type="submit">
				Process
			</button>
		</form>
</body>
</html>