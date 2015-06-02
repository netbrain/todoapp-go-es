angular.module('todoapp').factory('todoappws', ['$websocket',
    function($websocket) {
        var wsUrl = 'ws://'+location.hostname+(location.port ? ':'+location.port: '')+'/ws/'
        var datastream = $websocket(wsUrl)
        var collection = [];
        var eventMap = {}
        
        datastream.onMessage(function(message) {
            event = JSON.parse(message.data)
            for(eventType in eventMap){
                if(event.name === eventType){
                    for (var i = 0; i < eventMap[eventType].length; i++) {
                        eventMap[eventType][i](event)
                    };
                }
            }    
            collection.push(JSON.parse(message.data));
        });

        var changeSubsription = function(){
            datastream.send(JSON.stringify({
                events: Object.keys(eventMap),
            }));            
        }
        var methods = {
            on: function(event, fn){
                if(!eventMap[event]){                    
                    eventMap[event] = []
                    changeSubsription()
                }
                eventMap[event].push(fn)
            }
        };

        return methods
    }
])