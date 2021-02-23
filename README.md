# rvm-rails-assets

A Cloud Native Buildpack, compiling rails assets.
Requires "rvm-cnb" (Ruby language) and  "rvm-bundler-cnb" (Ruby bundler).

Initially based on a [paketo rails-assets](https://github.com/paketo-buildpacks/rails-assets) v0.1.0 (e4fe4db715e8dba19516a3cb72acb1963f8b36d2)

## Detect phase
- Project's "app/assets" directory must exist.
- Project's "Gemfile" must exist and contain gem "rails".
- Node-engine ("node" buildpack) is a dependency (Rails v5 mode)
- If project's "yarn.lock" exists, then Yarn ("yarn-install", "yarn" and "node") is a dependency instead of Node-engine above (Rails v6 mode).

## Build phase
- Preserves `RAILS_ENV` environment variable or sets it to `production` if not defined.
- Sets environment variable `RAILS_SERVE_STATIC_FILES=true`.
- Runs next command in RVM environment.
```shell
"bundle exec rails assets:precompile assets:clean"
```

 
## Buildpack building command example
*package.sh script downloads "jam" and "pack" binaries and uses them.*

```sh
$ scripts/package.sh -v "0.0.1"
```
It packs .tgz and .cnb to a "build" directory.

If only binaries are sufficient, then next command just compiles them to a "bin" directory:
```sh
$ scripts/build.sh
```

## Order group for `builder.toml` example
[[order]]

  [[order.group]]
  id = "com.anynines.buildpacks.gitcredentials"
  optional = true

  [[order.group]]
  id = "com.anynines.buildpacks.rvm"

  [[order.group]]
  id = "com.anynines.buildpacks.rvm-bundler"

  [[order.group]]
  id = "paketo-buildpacks/node-engine"
  optional = true

  [[order.group]]
  id = "paketo-buildpacks/yarn"
  optional = true

  [[order.group]]
  id = "paketo-buildpacks/yarn-install"
  optional = true

  [[order.group]]
  id = "com.anynines.buildpacks.rvm-rails-assets"
