var apiai = require('apiai');
// to asynchronously spawn a child process to speak back to the user
var spawn = require('child_process').spawn;
 
var app = apiai("b00a32a7fe1146789413024c1d42f027");

const ROOT_DIR = __dirname
const Sonus = require('sonus')

const speech = require('@google-cloud/speech')({
  projectId: 'streaming-speech-sample',
  keyFilename: './RonBot-731caf2fc154.json'
})

const hotwords = [{ file: 'resources/ron.pmdl', hotword: 'ron' }]
const language = "en-US"

//recordProgram can also be 'arecord' which works much better on the Pi and low power devices
const sonus = Sonus.init({ hotwords, language, recordProgram: "arecord" }, speech)

Sonus.start(sonus)
console.log('Say "' + hotwords[0].hotword + '"...')

sonus.on('hotword', (index, keyword) => spawn('aplay', ['ding.wav']))

sonus.on('partial-result', result => console.log("Partial", result))

sonus.on('final-result', result => {
	console.log("Final", result)
	var request = app.textRequest(result, {
	    sessionId: '1'
	});
	// sends request to our api.ai chat bot that responds with a relevant thing to be spoken
	request.on('response', function(response) {
		var speech = response.result.fulfillment.speech;
		    if (speech) {
		    	console.log("speaking...");
		        const say = spawn('./speech.sh', [speech])	
                        console.log(speech);
		    }
	});
	request.on('error', function(error) {
	// run /usr/bin/say "Sorry, something went wrong..."
    	const say = spawn('./speech.sh', ["Sorry, something went wrong..."]);
	    console.log(error);
	});

	if (result.includes("stop")) {
		Sonus.stop()
	}
	
	request.end();
})

