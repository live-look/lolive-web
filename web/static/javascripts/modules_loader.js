(function() {
  function _loadMainModules(eventName) {
    var moduleLoader = function () {
      for (var item in app.modules) {
        if (!app.modules[item][eventName]) { continue; }

        try { app.modules[item][eventName](); }
        catch(error) { console.trace(error); }
      }
    };

    setTimeout(moduleLoader, 10);
  }

  function _onWindowLoad() {
    var promise = new Promise((resolve, reject) => {
      setTimeout(function() {
        resolve();
      }, 1000);
    });

    promise.
      then(() => _loadMainModules('load')).
      catch((e) => console.error(e));
  }

  function _listener() {
    document.addEventListener('DOMContentLoaded', function(event) {
      _loadMainModules('ready');
    });
    window.addEventListener('load', _onWindowLoad);
  }

  app.modules && _listener();
})();
