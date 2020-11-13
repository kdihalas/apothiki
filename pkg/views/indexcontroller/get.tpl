<div class="columns is-centered">
  <div class="column is-two-thirds">
    <nav class="panel">
      <p class="panel-heading">
        Repositories
      </p>
      {{range $key, $val := .repos}}
      <a class="panel-block is-active" href="/{{$val}}/">
        <span class="panel-icon">
          <i class="fas fa-book" aria-hidden="true"></i>
        </span>
        {{$val}}
      </a>
      {{end}}
    </nav>
  </div>
</div>