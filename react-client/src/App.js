import * as React from 'react';
import './App.css';
import * as ReactDOM from 'react-dom';
import { Chat } from '@progress/kendo-react-conversational-ui';
import {sendToServer, checkInput} from './interact'

const user = {
  id: 1,
  avatarUrl: "https://cdn-icons-png.flaticon.com/512/147/147144.png"
};
const bot = {
  id: 0
};
const initialMessages = [{
  author: bot,
  suggestedActions: [{
    type: 'reply',
    value: 'ls /'
  }],
  timestamp: new Date(),
  text: "Welcome to SFS - Simple File System"
}];


function App() {
  const [messages, setMessages] = React.useState(initialMessages);

  const addNewMessage = async (event) => {
    let botResponse = Object.assign({}, event.message);
    let input = event.message.text;
    if (checkInput(input) == false){
      botResponse.text = "invalid syntax"
    } else {
      botResponse.text = await sendToServer(event.message.text);
    }
    botResponse.author = bot;
    setMessages([...messages, event.message]);
    setTimeout(() => {
      setMessages(oldMessages => [...oldMessages, botResponse]);
    }, 1000);
  };

  return <div>
        <Chat user={user} messages={messages} onMessageSend={addNewMessage} placeholder={"Type a command..."} width={600} />
      </div>;
}

export default App