# Tic-tac-toe

Backend server for a tic-tac-toe game. frontend is a simple .md file 😎 . 
Checkout the frontend repo [here](https://github.com/sridhar-sp/tic-tac-toe)


# Interested in playing Tic-tac-toe from Github Readme.md file. 

Checkout: https://github.com/sridhar-sp/tic-tac-toe

This is a learning attempt to see how interactive we can make the Github readme.md file.

Checkout the [medium article](https://betterprogramming.pub/play-tic-tac-toe-from-github-readme-md-file-754539603380)



# Backstory

Whenever I add an image link to the Github readme file it’s always getting replaced with ‘https://camo.githubusercontent.com/some-hash’ proxy url.

I wondered why it was happening like that. I mean, it’s okay to proxy an external url, but why are even static images placed in repositories getting a proxy url when they're referenced from a readme file? What is Github trying to achieve here?

Then I thought what would happen if there was no proxy server involved while serving the images. If there is no proxy, then the http request to fetch the image will directly come to our server(where we host the image) and we can read the http request and try to get a user ip address, which can be used for any tracking purposes and Github will have no control over it.

So I kind of understood the reason behind why Github added the image proxy server. At this stage, I wanted to know whether it’s possible to get any unique identification from the http request made from the user machine.

I hoped there would be a way to uniquely identify each session/machine. So I thought of building a simple game that can be played from within the markdown file itself to test that. The idea was to build a multiplayer game where each player would have their own state saved against a unique id. "Tic-tac-toe" game seems to be a good fit for this, since it's a well-known game, and it’s fairly straightforward to implement.

After I started the development, I quickly found out there was no way we could get a unique id from the http request fired from the readme markdown file. This is because each time a request comes from a random proxy server (when "Cache-Control" is set as "no-cache,no-store,must-revalidate"), and on top of that, cookies are not allowed either. So, without a unique id, there won’t be any states saved for each player, and individual game play is impossible. Therefore, only a single game play will be shared with the entire internet. (In other words, anyone in the world can see and play the game, but everyone will only see one game play; changes made by one will be reflected to everyone)

At this stage, I thought of abandoning the quest, but playing a game from the Github README.md file itself seems like a cool idea, even when the gameplay is shared with the entire internet. So I ran this idea by my friends, and they seemed interested in seeing how this would turn out, so I thought of investing some time in developing this.



<br/>

# Open source library used
```
* https://github.com/ajstarks/svgo
```
