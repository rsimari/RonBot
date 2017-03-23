Ron

John Joyce & Rob Simari

For project 2, we set out to build a simple interactive hardware bot, not unlike Amazon Echo. The philosophy behind it was to create a hardware assistant that does virtually all of its logic in the cloud, using an extremely simple client-server communication model. After much tinkering, we have finally built an MVP that will be vastly expanded upon and improved (maybe even for project 4?). In its current state, our bot runs on a Rasberry Pi wired up to an external microphone and speaker, which provide the user with an interface by which to interact with Ron. The functionality that is supported includes the ability to query for current weather anywhere in the world, to schedule reminder texts to be sent to a predefined phone number, the reading of reddit top posts, the current headlines and article descriptions for a multitude of common news outlets, and the ability to recite simple P. Bui wisecracks.

One may ask, how did we do all of this? The answer was APIs. Our bot relies on a slew of APIs to retrieve its data as well as to parse the input speach text into useful data that we are able to analyze ourselves. The entire process begins with the hardware itself. On the hardware runs a Node.js program (listen.js) that uses a library called 'Sonus' to listen for keywords (Ron, in this case) and record input. This is then forwarded on to the Google Speech API, which translates the audio into text. From here, the text is forwarded to the popular AI/chatbot service named API.ai. At this stage, the text is parsed, with the intent and parameters being extracted as per our desires. To achieve the desired functionality, we had to configure a Bot on API.ai's platform as well as provide it with webhooks to which it could place processed results. API.ai communicates its results to our Go webserver, residing on an EC2 instance. Each query hits hits just a single REST endpoint on our server at 'api/speech/'. Once we have received the POST from API.ai, we parse it by the desired action. These actions includes things such as 'reddit_top_post' & 'set_reminder'. From here, we decide how to proceed to perform the desired action. The fulfill these requests we tap into different APIs. For reddit, we obviously use the reddit api. For news, we use a google news aggregator api, for weather, we use a global weather api, etc. The most interesting of these is probably the logic required for the 'set_reminder' action. Requesting this employs 2 platforms - Iron.io & Twilio. Iron.io is essentially a worker queue API, providing a place to put pieces of code that can be remotely executed via API at some determined time. In our case, it allows us to call 'Go' code that hits the Twilio API at scheduled times in order to send reminder text messages. 

The unique feature of 'Go' that we chose to employ were Go rountines. In the context of the app, we use them when altering the file related to a user's preferences. In a multi-user scenario, this would keep the file I/O from blocking the webserver, which we thought was a valuable place to start with Go routines. In implementing this functionality, we were able to gain valuable experience with them. We'd like to use them elsewhere, particularly for handling requests in the context of a multi-user application.  

The greatest source of challenge behind the project were learning the various platforms we used (specifically API.ai, Iron.io, Twilio) and integrating them. There is still much work left to do on the project, and we plan to continue with it. The next steps would be to successfully enable authentication (partially implemented currently), generalize skills and create a micro package manager to enable such skills so that third party developers could independently create skills, to build a nice outer shell for the hardware/upgrade hardware, and to have multi-user support for the bot. As of now, it is highly personal, only supporting a single user. Finally, we'd like to create a paired iPhone application that could control user preferences/enable skills. 

How to run bot:

1. Get rasberry pi + mic + speaker
2. Run our 'audio' folder code on rasberry pi
3. Say 'Hey Ron'

How to run server code: 

1. Enter our 'GoServer' directory
2. Type 'Go build'
3. Run ./GoServer


