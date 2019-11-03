FROM gliderlabs/alpine

RUN apk add --no-cache --update wget build-base perl tzdata \
    postgresql-dev postgresql-client libpq perl-dev

RUN cpan App::cpanminus && \
  cpanm --notest -q App::Sqitch && \
  cpanm --notest -q DBD::Pg && \
  cpanm --notest -q TAP::Parser::SourceHandler::pgTAP

RUN apk del --purge postgresql-dev perl-dev build-base wget && \
  rm -rf ~/.cpan ~/.cpanm && \
  ln -sf /usr/share/zoneinfo/UTC  /etc/localtime

WORKDIR /src

