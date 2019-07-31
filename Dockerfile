FROM heroku/heroku:18-build as build

LABEL maintainer="Harry Brown <harrybrown98@gmail.com"

COPY . /app
WORKDIR /app

RUN mkdir -p /tmp/buildpack/heroku/go /tmp/build_cache /tmp/env
RUN curl https://codon-buildpacks.s3.amazonaws.com/buildpacks/heroku/go.tgz | tar xz -C /tmp/buildpack/heroku/go

RUN STACK=heroku-18 /tmp/buildpack/heroku/go/bin/compile /app /tmp/build_cache /tmp/env

FROM heroku/heroku:18

COPY --from=build /app /app
ENV HOME /app
WORKDIR /app
RUN useradd -m heroku
USER heroku

CMD /app/bin/harrybrown.com