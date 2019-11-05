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

      console.log("Local session description: ");
      console.log(btoa(JSON.stringify(peerConnection.localDescription)));

      //_enableButton();
    }
  }

  function _setupVideo() {
  }

  function _startSession() {
    $.ajax({
      method: 'POST',
      url: '/broadcasts',
      data: JSON.stringify({local_sdp: btoa(JSON.stringify(peerConnection.localDescription))}),
      contentType: 'application/json',
      dataType: 'json',
      success: function(data) {
        console.log('Success');
        console.log(data);

        try {
          peerConnection.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(data.remote_sdp))))
        } catch (e) {
          console.error(e);
        }
      }
    });
  }

  function _stopSession() {

  }

  self.load = function() {
    _createSession();
    _setupVideo();
  };

  return self;
})(app.modules.viewer || {});
