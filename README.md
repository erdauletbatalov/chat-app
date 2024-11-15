# chat-app
A chat room where users can send and receive messages in real time.

## Implementation of technologies:

* WebSockets: For instant message delivery.

* SSE: For displaying notifications about new users or events (e.g. "User joined").

* Long Polling: As a fallback option for devices or browsers without WebSocket support.

* Short Polling: For periodic checking of the number of unread messages.

## Additional: 

* Implement authorization, message history and support for multiple rooms.