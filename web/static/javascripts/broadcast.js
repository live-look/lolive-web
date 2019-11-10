app.modules.broadcast = (function(self) {
  var peerConnection = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  }),
  yourSelfVideo = document.getElementById('video'),
  broadcastingEnabled = false;

  function _createSession() {
    peerConnection.oniceconnectionstatechange = function(event) {
      console.log(peerConnection.iceConnectionState);
    }

    peerConnection.onicecandidate = function(event) {
      console.log(event);
      if (event.candidate === null) {
        _enableButton();
      }
    }
  }

  function _setupWebcam() {
    navigator.mediaDevices.getUserMedia({video: true, audio: false}).then(function(stream) {
      yourSelfVideo.srcObject = stream;

      peerConnection.addStream(stream);
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

  function _startSession() {
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

  function _stopSession() {

  }

  function _enableButton() {
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

  self.ready = function() {
    _createSession();
    _setupWebcam();
  };

  return self;
})(app.modules.broadcast || {});
