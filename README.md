# Tic-tac-toe

Backend server for a tic-tac-toe game. frontend is a simple .md file 😎 . 
Checkout the frontend repo [here](https://github.com/sridhar-sp/tic-tac-toe)


# Interested in playing Tic-tac-toe from Github Readme.md file. 

Checkout: https://github.com/sridhar-sp/tic-tac-toe

This is a learning attempt to see how interactive we can make the Github readme.md file.

# Backstory

Whenever I add an image link to the Github readme file it’s always getting replaced with ‘https://camo.githubusercontent.com/some-hash’ proxy url.

I wondered why it's happening like that. I mean it’s okay to proxy external urls but why even static images placed in repositories are getting a proxy url when it's referenced from readme file. What is Github trying to achieve here?

Then I thought what would happen if there is no proxy server involved while serving the images. If there is no proxy then the http request to fetch the image will directly come to our server and we can read the http request and try to get a user ip address which can be used for any tracking purposes and Github would have no control over it.

So I kind of understood the reason behind why Github added the image proxy server. At this stage I wanted to know whether it’s possible to get any unique identification from the http request made from the user machine.

I hoped there would be a way to uniquely identify each session/machine so I thought of building a simple game that can be played from the markdown file itself. Each player will have their own state saved against the unique id. For this “tic-tac-toe” seems to be a good fit, since everyone is aware of the game, and it's fairly straightforward to implement.

After I started the development I quickly found out there is no way we can get a unique id from the http request. This is because each time a request comes from a random proxy server and on top of that cookies are not allowed too. So without a unique id there won’t be any state saved for each player and individual game play is impossible, thereby a single game play will be shared with the entire internet.

At this stage I thought of abandoning the quest, but playing a game from the Github Readme.md file seems like a cool idea. So I ran this idea over with my friends and they seem to be interested to see how this turns out, so thought of investing some time to develop this.



<br/>

# Open source library used
```
* https://github.com/ajstarks/svgo
```
