* Drawing a fancy Chrome timeline.
  * Use the data format from Network.Network.requestWillBeSent and Network.responseReceived
    * begin: Network.requestWillBeSent
    * finish: Network.responseReceived's timing (ResourceTiming: dns and other information)
  * Visualization
    * resource type --> specific color
    * start and end as a bar with an assigned color
    * export as PNG or SVG file 
    * as cli tool and server service

* Related to
  * https://chromedevtools.github.io/devtools-protocol/tot/Network#type-ResourceTiming

