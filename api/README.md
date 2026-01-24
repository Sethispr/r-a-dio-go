This package handles all the comms w/ r-a-d.io. All it does is act like a "normal" browser" so it doesnt get blocked

- **Check Status:** Sees what song is playing and if the current dj is a person or hanyuu clanker

- **Search:** Finds songs in the api db from input

1. Finds csrf token from the html
2. Uses a proxy so you can cycle through every request with diff identity
3. Sends a request for your song to be played (you can also use the cart feature)
