# go-corona

## Introduction: 
A CLI to get an email containing the metrics associated with coronavirus of any country you specify<br>

------

## Pre-requisites:
Please have "go" installed in your computer. That's it!

------

## How to set it up:
 - Download the github repo with the following command:<br>
 ```git clone https://www.github.com/yashvardhan-kukreja/go-corona```
 - Then, run the following command:<br>
 ```bash commander.sh```

 - And, that's it !

-----

## How to run it:
 
 The command expects four parameter:
 - **email** -> the email id at which you want to receive alerts
 - **password** -> the password of that email account. <b>I promise I am not going to steal it. Look at the source code if you want :)</b>
 - **country** (default: India)-> The country for which you want to receive COVID alerts
 - **timeInSeconds** (default: 300 seconds) -> rate at which you want to receive covid alerts/email

 Example command:
 - Here, the email is "yash.98@gmail.com" and the password is "hello1234". And, let's say I want to get the alerts every 10 minutes (600 seconds) regarding the status of Covid in Canada.<br>
 So, this would be the command to do so:<br>
 ```go-corona --email yash.98@gmail.com --password hello1234 --country Canada --timeInSeconds 600 ```

-----

## Upcoming features:
 - Getting emails continuously at a defined rate
 - Mail delivery to multiple mail recipients

-----