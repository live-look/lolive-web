<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link rel="shortcut icon" href="https://static.camforchat.docker/assets/favicon.ico" type="image/x-icon">

    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <link rel="stylesheet" href="https://static.camforchat.docker/assets/stylesheets/main.css">

    <title>Camforchat - free webcams</title>
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <a class="navbar-brand">Camforchat</a>
      <ul class="navbar-nav">
      </ul>
      <ul class="navbar-nav ml-md-auto">
        <li class="nav-item">
          <a class="nav-link" href="/auth/register" rel="noopener">Register</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" href="/auth/login" rel="noopener">Login</a>
        </li>
      </ul>
    </nav>

    <div class="ab-forms">
      {{block "authboss" .}}{{end}}
    </div>
  </body>
</html>
