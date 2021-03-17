# rvm-rails-assets

A Cloud Native Buildpack, compiling rails assets.
Requires "rvm-cnb" (Ruby language) and  "rvm-bundler-cnb" (Ruby bundler).

Based on a [paketo rails-assets](https://github.com/paketo-buildpacks/rails-assets) v0.2.0 (6a2741cee08828ab718d13d53046ee0de773ac31)

## Detect phase
- Project's assets must be in any directory of `app/assets`, `app/javascript`, `lib/assets`, or `vendor/assets`.
- Project's "Gemfile" must exist and contain gem "rails".
- Node-engine is a dependency.
- If project's "yarn.lock" exists, then Yarn installs as a dependency ans `yarn install` runs.

## Build phase
- Preserves `RAILS_ENV` environment variable or sets it to `production` if not defined.
- Preserves `RAILS_SERVE_STATIC_FILES` environment variable or sets it to `true` if not defined.
- Runs next command in RVM environment.
```shell
"bundle exec rake assets:precompile assets:clean"
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
