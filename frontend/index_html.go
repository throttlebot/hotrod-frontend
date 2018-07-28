package frontend

import (
	"html/template"
)

var indexHTML = template.Must(template.New("index").Parse(`<html>
  <meta charset="ISO-8859-1">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0">
  <head>
    <title>HotROD - Rides On Demand</title>
    <script src="https://code.jquery.com/jquery-3.1.1.min.js"></script>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>

    <style>
.uuid { margin-top: 15px; }
.hotrod-button { padding: 20px; cursor: pointer; margin-top: 10px; }
.hotrod-button:hover { cursor: pointer; filter: brightness(85%); }
#hotrod-log { margin-top: 15px; }
#tip { margin-top: 15px; }
    </style>

  </head>
  <body>
    <div class="container">
      <div class="uuid alert alert-info"></div>
      <center>
        <h1>Hot R.O.D.</h1>
        <h4><em>Rides by Will</em></h4>
        <div class="row">
			{{range .Customers}}
				<div class="col-md-3 col-sm-6">
					<span
						class="btn btn-info btn-block hotrod-button"
						data-customer="{{.ID}}">{{.Name}}</span>
				</div>
			{{end}}
        </div>
        <div id="tip">Click on customer name above to order a car.</div>
        <div id="hotrod-log" class="lead"></div>
      </center>
    </div>
  </body>

  <script>

function formatDuration(duration) {
  var d = duration / (1000000 * 1000 * 60);
  var units = 'min';
  return Math.round(d) + units;
}

function Refund(request) {
    req = requestsCancellable[request]
    if (req == undefined) {
       $("#hotrod-log").prepend('<div class="fresh-car">Request ' + request + ' has already been refunded.</div>')
       return
    }
    var textBox = $($("#hotrod-log").prepend('<div class="fresh-car"><b>Canceling request ' + request + '</b></div>').children()[0])
	$.ajax('/refund?driver=' + req.Driver + '&customer=' + req.Customer, {
        method: 'GET',
        success: function(data, textStatus) {
          delete requestsCancellable[request]
	       textBox.html('Request '  + request + ' cancelled')
        },
        error: function(xhr, error){
            if (xhr.status != 200) {
         		textBox.html('Request '  + request + ' cancellation failed')
			} else {
                delete requestsCancellable[request]
				textBox.html('Request '  + request + ' cancelled')
			}
		},
    })
}

var clientUUID = Math.round(Math.random() * 10000);
var lastRequestID = 0;
var requestsCancellable = {}

$(".uuid").html("Your web client's id: <strong>" + clientUUID + "</strong>");

$(".hotrod-button").click(function(evt) {
  lastRequestID++;
  var requestID = clientUUID + "-" + lastRequestID;
  var freshCar = $($("#hotrod-log").prepend('<div class="fresh-car"><em>Dispatching a car...[req: '+requestID+']</em></div>').children()[0]);
  var customer = evt.target.dataset.customer;
  headers = {
      'jaeger-baggage': 'session=' + clientUUID + ', request=' + requestID
  };
  console.log(headers);
  var before = Date.now();
  $.ajax('/dispatch?customer=' + customer + '&nonse=' + Math.random(), {
    headers: headers,
    method: 'GET',
    success: function(data, textStatus) {
      var after = Date.now();
      console.log(data);
      requestsCancellable[" " + requestID] = {Driver: data.Driver, Customer: customer}
      var duration = formatDuration(data.ETA);
      freshCar.html('HotROD <b>' + data.Driver + '</b> arriving in ' + duration + ' [req: ' + requestID + ', latency: ' + (after-before) + 'ms] <a href="#" onclick="Refund(&#39 ' + requestID + '&#39); return false">Cancel</a>');
    },
  });
});

  </script>

</html>`))
