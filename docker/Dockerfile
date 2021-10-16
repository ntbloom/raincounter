# run raincounter from docker
FROM fedora:latest
RUN dnf upgrade -y
ENV TZ=America/New_York

COPY  ./raincounter /bin/raincounter
COPY ./configs/* /etc/rainbase/
