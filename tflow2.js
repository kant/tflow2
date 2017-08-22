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
const OpEqual = 0;
const OpUnequal = 1;
const OpSmaller = 2;
const OpGreater = 3;

var bdfields = [
        "SrcAddr", "DstAddr", "Protocol", "IntIn", "IntOut", "NextHop", "SrcAsn", "DstAsn",
        "NextHopAsn", "SrcPfx", "DstPfx", "SrcPort", "DstPort" ];

var fields = ["Router"].concat(bdfields)

function drawChart() {
    var query = $("#query").val();
    if (query == "" || query == "{}") {
        return;
    }

    var url = "/query?q=" + encodeURI(query)
    console.log(url);

    $.ajax({
        type: "GET",
        url: url,
        dataType: "text",
        success: function(rdata) {
            rdata = rdata.trim()
            pres = Papa.parse(rdata)

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

            var chart = new google.visualization.AreaChart(document.getElementById('chart_div'));
            chart.draw(data, options);
        }
    });
}

function populateForm() {
    var q = $("#query").val();
    if (q == "" || q == "{}") {
        return;
    }

    q = JSON.parse(q);
    $("#topn").val(q.TopN);
    for (var c in q['Cond']) {
        var fieldName = q['Cond'][c]['Field'];
        var operand = q['Cond'][c]['Operand'];
        if (fieldName == "Router") {
            operand = getRouterById(operand);
            if (operand == null) {
                return;
            }
        } else if (fieldName == "IntIn" || fieldName == "IntOut") {
            operand = getInterfaceById($("#Router").val(), operand);
            if (operand == null) {
                return;
            }
        } else if (fieldName == "Protocol") {
            operand = protocols[operand];
            if (operand == null) {
                return;
            }
        }

        $("#" + fieldName).val(operand);
    }
    loadInterfaceOptions();

    for (var f in q['Breakdown']) {
        $("#bd" + f).prop( "checked", true );
    }
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

$(document).ready(function() {
    var start = new Date(((new Date() / 1000) - 900)* 1000).toISOString().substr(0, 16)
    if ($("#TimeStart").val() == "") {
        $("#TimeStart").val(start);
    }

    var end = new Date().toISOString().substr(0, 16)
    if ($("#TimeEnd").val() == "") {
        $("#TimeEnd").val(end);
    }

    $.when(loadRouters(), loadProtocols()).done(function() {
        $("#Router").on('input', function() {
            loadInterfaceOptions();
        })
        populateForm();
    })

    $("#submit").on('click', submitQuery);

    google.charts.load('current', {
        'packages': ['corechart']
    });
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
    var query = {
        Cond: [],
        Breakdown: {},
        TopN: parseInt($("#topn").val())
    };

    console.log($("#TimeStart").val());
    var start = new Date($("#TimeStart").val());
    var end = new Date($("#TimeEnd").val());
    start = Math.round(start.getTime() / 1000);
    end = Math.round(end.getTime() / 1000);
    query['Cond'].push({
        Field: "Timestamp",
        Operator: OpGreater,
        Operand: start + ""
    });
    query['Cond'].push({
        Field: "Timestamp",
        Operator: OpSmaller,
        Operand: end + ""
    });

    for (var k in fields) {
        field = fields[k]

        tmp = $("#" + field).val();
        if (tmp == "") {
            continue;
        }
        if (field == "Router") {
            tmp = rtrs[tmp]['id'];
        } else if (field == "IntIn" || field == "IntOut") {
            tmp = getIntId($("#Router").val(), tmp)
            if (tmp == null) {
                return;
            }
        } else if (field == "Protocol") {
            tmp = getProtocolId(tmp);
            if (tmp == null) {
                return;
            }
        }
        query['Cond'].push({
            Field: field,
            Operator: OpEqual,
            Operand: tmp + ""
        });
    }

    for (var i = 0; i < bdfields.length; i++) {
        if (!$("#bd" + bdfields[i]).prop('checked')) {
            continue;
        }
        query['Breakdown'][bdfields[i]] = true;
    }

    console.log(query);
    $("#query").val(JSON.stringify(query));
    $("#form").submit();
}