# rvm-rails-assets

A Cloud Native Buildpack, compiling rails assets.
Requires "rvm-cnb" (Ruby language) and  "rvm-bundler-cnb" (Ruby bundler).

Based on a [paketo rails-assets](https://github.com/paketo-buildpacks/rails-assets) v0.5.0(fb423968e88ce3947c422d12fcc7ebe52f803a6a) with golang libraries updates actual on 20.06.2022.

## Detect phase
- Project's assets must be in any directory of `app/assets`, `app/javascript`, `lib/assets`, or `vendor/assets`.
- Project's `Gemfile` must exist. There is no need to Gem `rails` be defined directly in `Gemfile`.
- Node-engine buildpack is a dependency.
- If project's "yarn.lock" exists, then Yarn installs as a dependency and `yarn install` runs.

## Build phase
- Preserves `RAILS_ENV` environment variable or sets it to `production` if not defined.
- Preserves `RAILS_SERVE_STATIC_FILES` environment variable or sets it to `true` if not defined.
- Runs next command in RVM environment.
```shell
"bundle exec rake assets:precompile assets:clean"
```

Buildpack may try to reuse cached layer with assets if SHA256 checksum of assets directories shown above wasn't changed since last run. This behaviour may be enabled by setting environment variable `RAILS_ASSETS_DISABLE_CACHING=FALSE`.
If set, `bundle exec rake assets:precompile assets:clean` will not run if there weren't any changes in these directories.

Any of environment variables `RAILS_ENV`, `DB_ADAPTER` or `SECRET_KEY_BASE` may be passed as arguments
directly to a rake process if project requires them. For example passing `--env RAILS_ENV=production`
to a `pack` command, causes assets compile command to be changed as follows:
```shell
"RAILS_ENV=production bundle exec rake assets:precompile assets:clean"
```
 
## Buildpack building command example
*package.sh script downloads "jam" and "pack" binaries and uses them.*

```sh
$ scripts/package.sh -v "1.2.3"
```
It packs .tgz and .cnb to a "build" directory.


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

## Logging Configurations

To configure the level of log output from the **buildpack itself**, set the
`$BP_LOG_LEVEL` environment variable at build time either directly (ex. `pack
build my-app --env BP_LOG_LEVEL=DEBUG`) or through a [`project.toml`
file](https://github.com/buildpacks/spec/blob/main/extensions/project-descriptor.md)
If no value is set, the default value of `INFO` will be used.

The options for this setting are:
- `INFO`: (Default) log information about the progress of the build process
- `DEBUG`: log debugging information about the progress of the build process

```shell
$BP_LOG_LEVEL="DEBUG"
```
