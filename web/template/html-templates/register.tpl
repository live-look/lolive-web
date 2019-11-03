<div class="container">
  <div class="row justify-content-center">
    <div class="col-md-8">
      <div class="card">
          <div class="card-header">Register</div>
          <div class="card-body">
            <form action="{{mountpathed "register"}}" method="post">
              <div class="form-group row">
                {{with .errors}}{{with (index . "")}}{{range .}}<span>{{.}}</span><br />{{end}}{{end}}{{end -}}
                <label for="name" class="col-md-4 col-form-label text-md-right">Name:</label>
                <div class="col-md-6">
                  <input name="name" type="text" class="form-control" value="{{with .preserve}}{{with .name}}{{.}}{{end}}{{end}}" />
                  {{with .errors}}{{range .name}}<span>{{.}}</span><br />{{end}}{{end -}}
                </div>
              </div>

              <div class="form-group row">
                <label for="email" class="col-md-4 col-form-label text-md-right">E-mail:</label>
                <div class="col-md-6">
                  <input name="email" type="text" class="form-control" value="{{with .preserve}}{{with .email}}{{.}}{{end}}{{end}}" />
                  {{with .errors}}{{range .email}}<span>{{.}}</span><br />{{end}}{{end -}}
                </div>
              </div>

              <div class="form-group row">
                <label for="password" class="col-md-4 col-form-label text-md-right">Password:</label>
                <div class="col-md-6">
                  <input name="password" type="password" class="form-control" />
                  {{with .errors}}{{range .password}}<span>{{.}}</span><br />{{end}}{{end -}}
                </div>
              </div>

              <div class="form-group row">
                <label for="confirm_password" class="col-md-4 col-form-label text-md-right">Confirm Password:</label>
                <div class="col-md-6">
                  <input name="confirm_password" type="password" class="form-control" />
                  {{with .errors}}{{range .confirm_password}}<span>{{.}}</span><br />{{end}}{{end -}}
                </div>
              </div>

              <div class="col-md-6 offset-md-4">
                <input type="submit" value="Register" class="btn btn-primary">
                <a href="/" class="btn btn-link">Cancel</a>
              </div>

                {{with .csrf_token}}<input type="hidden" name="csrf_token" value="{{.}}" />{{end}}
            </form>
          </div>
      </div>
    </div>
  </div>
</div>
