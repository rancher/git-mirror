FROM gcr.io/google-containers/fluentd-elasticsearch:v2.4.0

RUN apt-get update
RUN apt-get install -y build-essential libgeoip-dev ruby-dev libmaxminddb-dev libmaxminddb0 nfs-common
RUN apt-get upgrade -y

RUN apt-get purge -y --auto-remove  && clean-apt 

RUN fluent-gem install fluent-plugin-geoip
RUN fluent-gem install fluent-plugin-influxdb
