import React, { Component } from 'react';
import './App.css';
import RunnerTimes from './views/RunnerTimes'

class App extends Component {
  constructor(props) {
    super(props);

    this.runnerTimes = React.createRef();
    this.connect = this.connect.bind(this);
  }

  componentDidMount() {
    this.connect();
  }

  connect = cb => {
    var socket = new WebSocket("ws://127.0.0.1:10000/ws");

    socket.onopen = msg => {
      console.log("Successfully Connected");
    };

    socket.onmessage = msg => {
      this.runnerTimes.current.getRunnerTimes();
    };

    socket.onclose = event => {
      console.log("Socket Closed Connection: ", event);
    };

    socket.onerror = error => {
      console.log("Socket Error: ", error);
    };
  };

  render() {
    return (
      <div className="root">
        <RunnerTimes ref={this.runnerTimes} />
      </div>
    )
  }
}

export default App;