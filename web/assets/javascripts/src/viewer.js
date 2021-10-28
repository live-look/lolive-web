app.modules.viewer = (function(self) {
  var peerConnection = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  }),
  broadcastVideo = document.getElementById('video');

  function _createSession() {
    peerConnection.oniceconnectionstatechange = function(event) {
      console.log(peerConnection.iceConnectionState);
    }

    // here we receive ICE candidates
    peerConnection.onicecandidate = function(event) {
      console.log(event);
      if (event.candidate === null) {
        console.log(peerConnection.localDescription);

        _startSession();
      }
    }

    peerConnection.addTransceiver('video', {'direction': 'sendrecv'});
    peerConnection.ontrack = function(event) {
      console.log("On track");
      console.log(event);

      broadcastVideo.srcObject = event.streams[0];
      broadcastVideo.autoplay = true;
      broadcastVideo.controls = true;
    }

    peerConnection.createOffer().then(function(d) {
      peerConnection.setLocalDescription(d);
    }).
    catch(function(e) { console.error(e); });
  }

  function _startSession() {
    fetch('/broadcasts/' + app.config.viewer.currentBroadcastId + '/viewers', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      },
      body: JSON.stringify({local_sdp: btoa(JSON.stringify(peerConnection.localDescription))})
    }).
      then(response => response.json()).
      then(result => {
        console.log("set remote description");
        try {
          peerConnection.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(result.remote_sdp))))
        } catch (e) {
          console.error(e);
        }
      });
  }

  function _stopSession() {

  }

  self.ready = function() {
    console.log("ready");

    _createSession();
  };

  return self;
})(app.modules.broadcast || {});
