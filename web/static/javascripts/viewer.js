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

    peerConnection.onicecandidate = function(event) {
      if (event.candidate !== null) {
        return
      }
    }

    peerConnection.addTransceiver('video', {'direction': 'recvonly'});
    peerConnection.createOffer().
      then(function(d) {
        peerConnection.setLocalDescription(d);

        _startSession();
      }).
      catch(function(e) { console.error(e); });
    peerConnection.ontrack = function(event) {
      broadcastVideo.srcObject = event.streams[0];
      broadcastVideo.autoplay = true;
      broadcastVideo.controls = true;
    }
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
    _createSession();
  };

  return self;
})(app.modules.broadcast || {});
