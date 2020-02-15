'use strict';

const peerConnection = new RTCPeerConnection({
  iceServers: [
    {
      urls: 'stun:stun.l.google.com:19302'
    }
  ]
});

class Broadcast extends React.Component {
  constructor(props) {
    super(props);
    this.state = { broadcastingEnabled = false };
  }

  createSession() {
    peerConnection.oniceconnectionstatechange = function(event) {
      console.log(peerConnection.iceConnectionState);
    }

    peerConnection.onicecandidate = function(event) {
      if (event.candidate === null) {
        _enableButton();
      }
    }
  }

  setupWebcam() {
    let yourSelfVideo = document.getElementById('video');

    navigator.mediaDevices.getUserMedia({video: true, audio: false}).then(function(stream) {
      yourSelfVideo.srcObject = stream;

      // addStream is obsolete
      // peerConnection.addStream(stream);
      for (const track of stream.getTracks()) {
        peerConnection.addTrack(track);
      }

      peerConnection.createOffer().then(function(description) {
        peerConnection.setLocalDescription(description);
      }).catch(console.error);
    }).catch(function(message) {
      var alertDangerEl = document.getElementsByClassName('js-alert-danger')[0];

      if (alertDangerEl) {
        alertDangerEl.innerText = 'Sorry, but you have no any webcam. Plug device and try to reload page.';
        alertDangerEl.style.display = 'block';
      }

      console.error(message);
    });
  }

  startSession() {
    fetch('/broadcasts', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      },
      body: JSON.stringify({local_sdp: btoa(JSON.stringify(peerConnection.localDescription))})
    }).
      then(response => response.json()).
      then(result => {
        try {
          peerConnection.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(result.remote_sdp))))
        } catch (e) {
          console.error(e);
        }
      });
  }

  stopSession() {
    peerConnection.close();
  }

  enableButton() {
    var btn = document.getElementsByClassName('js-switch-broadcast')[0]

    btn.addEventListener('click', (event) => {
      if (broadcastingEnabled) {
        btn.innerText = 'Start broadcasting';
        btn.classList.add('btn-outline-success');
        btn.classList.remove('btn-outline-danger');

        broadcastingEnabled = false;

        _stopSession();
      } else {
        btn.innerText = 'Stop broadcasting';
        btn.classList.add('btn-outline-danger');
        btn.classList.remove('btn-outline-success');
        broadcastingEnabled = true;

        _startSession();
      }
    });

    btn.removeAttribute('disabled');
  }

  render() {
    return (
      <div>
        <p>Hello, broadcaster!</p>
      </div>
    )
  }
}
