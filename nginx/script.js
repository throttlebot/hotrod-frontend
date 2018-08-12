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
    $.ajax('/api/refund?driver=' + req.Driver + '&customer=' + req.Customer, {
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

$.ajax('/api/customers', {
    method: 'GET',
    dataType: 'json',
    success: function(data, textStatus) {
        for (var i = 0; i < data.length; i++) {
            $("#customer-btns").append('<div class="col-md-3 col-sm-6">' +
                '<span class="btn btn-info btn-block hotrod-button" data-customer="' + data[i]["ID"] + '">' + data[i]["Name"] +
                '</span></div>')
        };
    },
});


$(document).on('click', '.hotrod-button', function(evt) {
    console.log("Hello");
    lastRequestID++;
    var requestID = clientUUID + "-" + lastRequestID;
    var freshCar = $($("#hotrod-log").prepend('<div class="fresh-car"><em>Dispatching a car...[req: '+requestID+']</em></div>').children()[0]);
    var customer = evt.target.dataset.customer;
    headers = {
        'jaeger-baggage': 'session=' + clientUUID + ', request=' + requestID
    };
    console.log(headers);
    var before = Date.now();
    $.ajax('/api/dispatch?customer=' + customer + '&nonse=' + Math.random(), {
        headers: headers,
        method: 'GET',
        success: function(data, textStatus) {
            var after = Date.now();
            console.log(data);
            requestsCancellable[" " + requestID] = {Driver: data.Driver, Customer: customer}
            var duration = formatDuration(data.ETA);
            freshCar.html('HotROD <b>' + data.Driver + '</b> arriving in ' + duration + ' [req: ' + requestID + ', latency: ' + (after-before) + 'ms] <a href="#" onclick="Refund(&#39 ' + requestID + '&#39); return false">Cancel</a>');
            $("#map").attr("src","/map?t=" + new Date().getTime());
            driverX = data.DriverLocation.split(",")[0];
            driverY = data.DriverLocation.split(",")[1];
            customerX = data.CustomerLocation.split(",")[0];
            customerY = data.CustomerLocation.split(",")[1];
            display_route(driverX, driverY, customerX, customerY);
        },
    });
});



var margin = {top: 30, right: 30, bottom: 30, left: 50},
    width = 600 - margin.left - margin.right,
    height = 600 - margin.top - margin.bottom;

var x = d3.scale.linear().range([0, width]);
var y = d3.scale.linear().range([height, 0]);

var xAxis = d3.svg.axis().scale(x).orient("bottom").ticks(5);
var yAxis = d3.svg.axis().scale(y).orient("left").ticks(5);

var svg = new_svg();

function new_svg() {
    d3.select("#svg > *").remove();
    return d3.select("#svg")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");
}

function color_link(link) {
    if (link["color"] != undefined) {
        return "stroke:" + link.color +  ";stroke-width:3";
    }
    return "stroke:gray;stroke-width:1";
}

function RenderGraph(json) {

    var idToNode = {};
    json.nodes.forEach(function(node) {
        idToNode[node.id] = node;
    });

    x.domain([d3.min(json.nodes, function(d) { return d.x; } ), d3.max(json.nodes, function(d) { return d.x; } )]);
    y.domain([d3.min(json.nodes, function(d) { return d.y; } ), d3.max(json.nodes, function(d) { return d.y; } )]);

    svg.selectAll("circle").data(json.nodes)
        .enter().append("circle")
        .attr("r", 0.5)
        .attr("cx", function(d) { return x(d.x); })
        .attr("cy", function(d) { return y(d.y); })
        .attr("style", "fill:gray;");

    svg.selectAll("line").data(json.links)
        .enter().append("line")
        .attr("x1", function(d) { return x(idToNode[d.source].x); })
        .attr("y1", function(d) { return y(idToNode[d.source].y); })
        .attr("x2", function(d) { return x(idToNode[d.target].x); })
        .attr("y2", function(d) { return y(idToNode[d.target].y); })
        .attr("style", color_link);

    axisOffset = 10;

    svg.append("g")
        .attr("transform", "translate(0," + (height + axisOffset) + ")")
        .call(xAxis);

    svg.append("g")
        .attr("transform", "translate(-" + axisOffset + ",0)")
        .call(yAxis);
}

function RenderPath(json) {
    function keyfunc(source, target) {
        return source.toString() + "-" + target.toString();
    }

    var linkSet = new Set();
    var nodeSet = new Set();

    json.links.forEach(function(l) {
        linkSet.add(keyfunc(l.source, l.target))
    });

    json.nodes.forEach(function(n) {
        nodeSet.add(n.id)
    });

    function color_route_link(link) {

        if (linkSet.has(keyfunc(link.source, link.target)) || linkSet.has(keyfunc(link.target, link.source))) {
            return "stroke:blue;stroke-width:5";
        }
        return "stroke:gray;stroke-width:1";
    }

    svg.selectAll("circle")
        .attr("style", function(node) {
            if (nodeSet.has(node.id)) {
                return "fill:blue;"
            }
            return "fill:gray;"
        })

    svg.selectAll("line")
        .attr("style", color_route_link)

}

function display_route(x, y, z, w) {
    svg = new_svg();
    $.ajax('/map/location', {
        method: 'GET',
        dataType: 'json',
        data: {"x": x, "y": y, "z": z, "w": w},
        success: function(data, textStatus) {
            RenderGraph(data);
            $.ajax('/map/route', {
                method: 'GET',
                dataType: 'json',
                data: {"x": x, "y": y, "z": z, "w": w},
                success: function(data, textStatus) {
                    RenderPath(data);
                },
            });
        },
    });
}

function display_map(x, y, z, w) {
    svg = new_svg();
    $.ajax('/map/location', {
        method: 'GET',
        dataType: 'json',
        data: {"x": x, "y": y, "z": z, "w": w},
        success: function(data, textStatus) {
            RenderGraph(data);
        }
    });
}

display_map(0, 0, 5 , 5);