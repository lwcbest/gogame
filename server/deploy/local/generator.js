//node generator.js
var env = 'local';
var exefile = '../../server';
var confs = require('../../conf/'+env+'/clusterconf.json');
var fs = require('fs');

var _RESULT = {};

var masterConf = confs["master"];
var servers = confs["servers"];

var process = {
  "name"              : "?",
  "cwd"               : "./",
  "script"            : "",
  "instances"         : 1,
  "exec_mode"         : "fork",
  "interpreter"       :exefile,
  "watch"             : false,
  "env"               : {
    "name": "?"
  },
  "merge_logs"        : true,
  "autorestart"       : true,
  "max_memory_restart": "1G"
};
var destHost = masterConf.host;
process.name = "master";
process.env.name = "master";
process.script = "./"+destHost+".json";

if (_RESULT.hasOwnProperty(destHost)) {
  _RESULT[destHost].apps.push(process);
} else {
  _RESULT[destHost] = {};
  _RESULT[destHost]['apps'] = [];
  _RESULT[destHost].apps.push(process);
}

for (var serverName in servers) {
  console.log(serverName)
  var process = {
    "name"              : "?",
    "cwd"               : "./",
    "script"            : "",
    "instances"         : 1,
    "exec_mode"         : "fork",
    "interpreter"       :exefile,
    "watch"             : false,
    "env"               : {
      "name": "?"
    },
    "merge_logs"        : true,
    "autorestart"       : true,
    "max_memory_restart": "1G"
  };
    var server = servers[serverName];
    destHost = server.host;
    process.name = String(serverName);
    process.env.name = String(serverName);
    process.script = "./"+destHost+".json";
    if (_RESULT.hasOwnProperty(destHost)) {
      _RESULT[destHost].apps.push(process);
    } else {
      _RESULT[destHost] = {};
      _RESULT[destHost]['apps'] = [];
      _RESULT[destHost].apps.push(process);
    }
 };


for (var host in _RESULT) {
  fs.writeFileSync('./' + host + '.json', JSON.stringify(_RESULT[host], null, 2), {
    flag: 'w'
  });
}