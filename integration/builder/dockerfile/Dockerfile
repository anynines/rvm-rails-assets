FROM ubuntu:bionic

ARG package_args='--allow-downgrades --allow-remove-essential --allow-change-held-packages --no-install-recommends'

RUN echo "debconf debconf/frontend select noninteractive" | debconf-set-selections && \
  export DEBIAN_FRONTEND=noninteractive && \
  apt-get -y $package_args update && \
  apt-get -y $package_args install locales && \
  locale-gen en_US.UTF-8 && \
  update-locale LANG=en_US.UTF-8 LANGUAGE=en_US.UTF-8 LC_ALL=en_US.UTF-8 && \
  apt-get -y $package_args install \
  autoconf \
  automake \
  bison \
  build-essential \
  ca-certificates \
  curl \
  gawk \
  git \
  graphicsmagick-libmagick-dev-compat \
  jq \
  libgmp-dev \
  libgdbm-dev \
  libffi-dev \
  libmagickwand-dev \
  libmysqlclient-dev \
  libncurses5-dev \
  libpq-dev \
  libreadline7 \
  libreadline-dev \
  libssl1.1 \
  libssl1.0.0 \
  libsqlite3-dev \
  libssl-dev \
  libxml2-dev \
  libxslt-dev \
  libyajl-dev \
  libyaml-0-2 \
  libyaml-dev \
  ruby-dev \
  sqlite3 \
  openssl \
  tzdata \
  xz-utils \
  zlib1g-dev && \
  rm -rf /var/lib/apt/lists/* /tmp/*

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install -y nodejs

RUN apt-get clean

RUN curl -sL -o /usr/local/bin/yj https://github.com/sclevine/yj/releases/latest/download/yj-linux \
  && chmod +x /usr/local/bin/yj
