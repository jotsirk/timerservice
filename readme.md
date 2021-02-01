The server file is at time_service folder. Run 'go run server.go' there and then navigate to client folder
and run 'npm start'.
To create random times for people who made it into the corridor. You can do that from '...\submitRunnerTime'
There a button 'add time' that creates a random person who made it into the finish corridor. To finish his time, click on the
'finish run' button on the grid.
client side datetime is unformatted at the moment, because i didnt want to put too much time into it. Same goes for the screen load everytime the socket sends a message. It was suppoused to send an updated object but unfortunately state changes were not recognized and it stayed this way.