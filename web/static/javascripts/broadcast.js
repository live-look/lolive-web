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
      if (event.candidate !== null) {
        return
      }

      console.log("Local session description: ");
      console.log(btoa(JSON.stringify(peerConnection.localDescription)));

      _enableButton();
    }
  }

  function _setupWebcam() {
    navigator.mediaDevices.getUserMedia({video: true, audio: true}).then(function(stream) {
      yourSelfVideo.srcObject = stream;

      peerConnection.addStream(stream);
      peerConnection.createOffer().then(function(description) {
        peerConnection.setLocalDescription(description);
      }).catch(console.error);
    }).catch(function(message) {
      var $alertDanger = $('.js-alert-danger');

      $alertDanger.text('Sorry, but you have no any webcam. Plug device and try to reload page.').show();

      console.error(message);
    });
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

  function _enableButton() {
    $(document).on('click', '.js-switch-broadcast', function(event) {
      var $this = $(this);

      if (broadcastingEnabled) {
        $this.text('Start broadcasting').addClass('btn-outline-success').removeClass('btn-outline-danger');
        broadcastingEnabled = false;

        _stopSession();
      } else {
        $this.text('Stop broadcasting').addClass('btn-outline-danger').removeClass('btn-outline-success');
        broadcastingEnabled = true;

        _startSession();
      }
    });

    $('.js-switch-broadcast').prop('disabled', false);
  }

  self.load = function() {
    _createSession();
    _setupWebcam();
  };

  return self;
})(app.modules.broadcast || {});
