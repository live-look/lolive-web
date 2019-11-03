<div class="container">
  <div class="row justify-content-center">
    <div class="col-md-8">
      <div class="card">
          <div class="card-header">Sign in</div>
          <div class="card-body">
            <form action="{{mountpathed "login"}}" method="post">
              {{with .error}}{{.}}<br />{{end}}
              <div class="form-group row">
                <label for="name" class="col-md-4 col-form-label text-md-right">Email:</label>
                <div class="col-md-6">
                  <input type="text" class="form-control" name="email" value="{{.primaryIDValue}}">
                </div>
              </div>

              <div class="form-group row">
                <label for="password" class="col-md-4 col-form-label text-md-right">Password:</label>
                <div class="col-md-6">
                  <input type="password" class="form-control" name="password">
                </div>
              </div>

              {{with .csrf_token}}<input type="hidden" name="csrf_token" value="{{.}}" />{{end}}
              {{with .modules}}{{with .remember}}<input type="checkbox" name="rm" value="true"> Remember Me</input><br />{{end}}{{end -}}
              {{with .redir}}<input type="hidden" name="redir" value="{{.}}" />{{end}}

              <div class="col-md-6 offset-md-4">
                <button type="submit" class="btn btn-primary">Login</button>
                {{with .modules}}{{with .recover}}<a href="{{mountpathed "recover"}}">Recover Account</a>{{end}}{{end -}}
                {{with .modules}}{{with .register}}<a href="{{mountpathed "register"}}" class="btn btn-link">Register Account</a>{{end}}{{end -}}
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

