(function () {

var espeak = require('espeak');
 
// optionally set the path to the `espeak` cli program if it's not in your PATH 
// espeak.cmd = '/usr/bin/espeak'; 
 
espeak.speak('hello world', function(err, wav) {
  if (err) return console.error(err);
  
  // get the raw binary wav data 
  var buffer = wav.buffer;
  
  // get a base64-encoded data URI 
  var dataUri = wav.toDataUri();
});
 
// optionally add custom cli arguments for things such as pitch, speed, wordgap, etc. 
espeak.speak('hello world, slower', ['-p 60', '-s 90', '-g 30'], function(err, wav) {});

})();
 