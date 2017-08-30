// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

var query;
var protocols;
var availableProtocols = [];
var rtrs;
var routers = [];
var interfaces = [];

function drawChart() {
    var query = location.href.split("#")[1]
    if (!query) {
        return;
    }

    $.ajax({
        type: "GET",
        url: "/query?" + query,
        dataType: "text",
        success: function(rdata, status, xhr) {
            if (rdata == undefined) {
                $("#chart_div").text("No data found")
                return
            }
            renderChart(rdata)
        },
        error: function(xhr) {
            $("#chart_div").text(xhr.responseText)
        }
    })
}

function renderChart(rdata) {
    pres = Papa.parse(rdata.trim())

    var data = [];
    for (var i = 0; i < pres.data.length; i++) {
        for (var j = 0; j < pres.data[i].length; j++) {
            if (j == 0) {
                data[i] = [];
            }
            x = pres.data[i][j];
            if (i != 0) {
                if (j != 0) {
                    x = parseInt(x)
                }
            }
            data[i][j] = x;
        }
    }

    data = google.visualization.arrayToDataTable(data);

    var options = {
        isStacked: true,
        title: 'NetFlow bps of top flows',
        hAxis: {
            title: 'Time',
            titleTextStyle: {
                color: '#333'
            }
        },
        vAxis: {
            minValue: 0
        }
    };

    new google.visualization.AreaChart(document.getElementById('chart_div')).draw(data, options);
}

// source: https://stackoverflow.com/a/26849194
function parseParams(str) {
    return str.split('&').reduce(function (params, param) {
        var paramSplit = param.split('=').map(function (value) {
            return decodeURIComponent(value.replace('+', ' '));
        });
        params[paramSplit[0]] = paramSplit[1];
        return params;
    }, {});
}

function populateForm() {
    var query = location.href.split("#")[1]
    if (!query) {
        return;
    }

    var params = parseParams(query);
    
    for (var key in params) {
        var value = params[key]
        
        if (key.match(/^Timestamp/)){
            timezoneOffset = (new Date()).getTimezoneOffset()*60
            value = formatTimestamp(new Date((parseInt(value) - timezoneOffset )*1000))
        } else if (key == "Router") {
            value = getRouterById(value);
            if (value == null) {
                continue;
            }
        } else if (key == "IntIn" || key == "IntOut") {
            value = getInterfaceById($("#Router").val(), value);
            if (value == null) {
                continue;
            }
        } else if (key == "Protocol") {
            value = protocols[value];
            if (value == null) {
                continue;
            }
        } else if (key == "Breakdown") {
            var breakdown = value.split(",")
            for (var i in breakdown) {
                $("#bd"+breakdown[i]).attr("checked", true)
                continue
            }
        }

        $("#" + key.replace(".","_")).val(value);
    }
    loadInterfaceOptions();
}

function loadInterfaceOptions() {
    var rtr = $("#Router").val();
    interfaces = [];
    if (!rtrs[rtr]) {
        return;
    }
    for (var k in rtrs[rtr]["interfaces"]) {
        interfaces.push(rtrs[rtr]["interfaces"][k]);
    }

    $("#IntIn").autocomplete({
        source: interfaces
    });

    $("#IntOut").autocomplete({
        source: interfaces
    });
}

function loadProtocols() {
    return $.getJSON("/protocols", function(data) {
        protocols = data;
        for (var k in protocols) {
            availableProtocols.push(protocols[k]);
        }

        $("#Protocol").autocomplete({
            source: availableProtocols
        });
    });
}

function loadRouters() {
    return $.getJSON("/routers", function(data) {
        rtrs = data;
        for (var k in rtrs) {
            routers.push(k);
        }

        $("#Router").autocomplete({
            source: routers,
            change: function() {
                loadInterfaceOptions();
            }
        });
    });
}

function formatTimestamp(date) {
    return date.toISOString().substr(0, 16)
}

$(document).ready(function() {
    var start = formatTimestamp(new Date(((new Date() / 1000) - 900)* 1000))
    if ($("#Timestamp_gt").val() == "") {
        $("#Timestamp_gt").val(start);
    }

    var end = formatTimestamp(new Date())
    if ($("#Timestamp_lt").val() == "") {
        $("#Timestamp_lt").val(end);
    }

    $.when(loadRouters(), loadProtocols()).done(function() {
        $("#Router").on('input', function() {
            loadInterfaceOptions();
        })
        populateForm();
    })

    $("form").on('submit', submitQuery);

    google.charts.load('current', {
        'packages': ['corechart']
    });
    
    window.onhashchange = function () {
        populateForm()
        google.charts.setOnLoadCallback(drawChart);
    }

    google.charts.setOnLoadCallback(drawChart);
});

function getProtocolId(name) {
    for (var k in protocols) {
        if (protocols[k] == name) {
            return k;
        }
    }
    return null;
}

function getIntId(rtr, name) {
    if (!rtrs[rtr]) {
        return null;
    }
    for (var k in rtrs[rtr]['interfaces']) {
        if (rtrs[rtr]['interfaces'][k] == name) {
            return k;
        }
    }
    return null;
}

function getRouterById(id) {
    for (var k in rtrs) {
        if (rtrs[k]['id'] == id) {
            return k;
        }
    }
    return null;
}

function getInterfaceById(router, id) {
    return rtrs[router]['interfaces'][id];
}

function submitQuery() {
    var breakdown = []
    var query = {};

    $(".in input").each(function(){
        var field = this.id.replace("_",".")
        var value = this.value

        if (value == "") {
            return;
        }
        
        if (this.id.match(/^Timestamp/)){
            value = Math.round(new Date(value).getTime() / 1000)
        } else if (field == "Router") {
            value = rtrs[value]['id'];
        } else if (field == "IntIn" || field == "IntOut") {
            value = getIntId($("#Router").val(), value)
            if (value == null) {
                return;
            }
        } else if (field == "Protocol") {
            value = getProtocolId(value);
            if (value == null) {
                return;
            }
        }
        query[field] = value + ""
    })

    $(".bd input:checked").each(function(){
        breakdown.push(this.id.replace(/^bd/,""));
    })
    if (breakdown.length) {
        query.Breakdown = breakdown.join(",")
    }

    location.href = "#" + jQuery.param(query)
    return false
}