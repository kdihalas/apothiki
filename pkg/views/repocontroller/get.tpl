<div class="columns is-centered">
    <div class="column is-two-thirds">
        <nav class="panel">
            <p class="panel-heading">
                Tags for repository {{.repo}}
            </p>
            {{range $key, $val := .tags}}
                <a class="panel-block is-active">
                <span class="panel-icon">
                  <i class="fas fa-tag" aria-hidden="true"></i>
                </span>
                    {{$val}}
                </a>
            {{end}}
        </nav>
    </div>
</div>