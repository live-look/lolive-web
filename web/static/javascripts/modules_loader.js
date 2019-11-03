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
    var deferred = $.Deferred().resolve();

    deferred.always(_loadMainModules.bind(null, 'load'));

    setTimeout(function() {
      deferred.resolve();
    }, 1000);
  }

  function _listener() {
    $(document).ready(_loadMainModules.bind(null, 'ready'));
    $(window).on('load', _onWindowLoad);
  }

  app.modules && _listener();
})();
