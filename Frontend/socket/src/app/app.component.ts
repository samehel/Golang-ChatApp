import { Component, OnInit, OnDestroy } from '@angular/core';
import { SocketService } from './socket.service';
import { FormsModule } from '@angular/forms';
import { NgFor } from '@angular/common';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [FormsModule, NgFor],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})


// Creating the class to handle our frontend of the chat app
export class AppComponent implements OnInit, OnDestroy {
  
  public messages: Array<any>;
  public chatBox: string;

  // Initializing our messages array as an empty array and our chatbox as an empty string when the class is initialized
  public constructor(private socket: SocketService) {
    this.messages = [];
    this.chatBox = "";
  }

  // We will subscribe to the event listener and depending on the event we will carry out a specific action
  public ngOnInit(): void {
    this.socket.getEventListener().subscribe(event => {
      if(event.type == "message") {
        let data = event.data.content;
        if(event.data.sender)
          data = event.data.sender + ": " + data;
        this.messages.push(data);
      }
      if(event.type == "close")
        this.messages.push("User has disconnected from the channel");
      if(event.type == "open")
        this.messages.push("User has connected to the channel");
    });
  }

  // We will close the socket if the connection was closed for some reason
  public ngOnDestroy(): void {
    this.socket.close();
  }

  // We will send the message that the user has written in our chat box the empty it out
  public send(): void {
    if(this.chatBox) {
      this.socket.send(this.chatBox);
      this.chatBox = "";
    }
  }
  
  // We will check whether the message send was from the user or from our system, if its our system we will display it in bold
  public isSystemMessage(message: string): string {
    return message.startsWith("/") ? "<strong>" + message.substring(1) + "</strong>" : message;
  }

}
