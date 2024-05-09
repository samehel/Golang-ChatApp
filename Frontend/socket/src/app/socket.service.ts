import { Injectable, EventEmitter } from '@angular/core';

/* 
    Marking it as available to inject as a dependency, 
    based on the socket, a certain event is set to trigger
*/ 
@Injectable({
    providedIn: 'root'
})
export class SocketService {

    private socket: WebSocket;
    private listener: EventEmitter<any> = new EventEmitter();

    public constructor() {
        console.log("SocketServer Insttantiated")
        this.socket = new WebSocket("ws://localhost:12345/ws");
        this.socket.onopen = event => {
            console.log("Websocket connection established")
            this.listener.emit({"type": "open", "data": event});
        }
        this.socket.onclose = event => {
            this.listener.emit({"type": "close", "data": event});
        }
        this.socket.onmessage = event => {
            this.listener.emit({"type": "message", "data": JSON.parse(event.data)});
        }
    }

    public send(data: string): void {
        this.socket.send(data);
    }

    public close(): void {
        this.socket.close();
    }

    public getEventListener(): EventEmitter<any> {
        return this.listener;
    }

}