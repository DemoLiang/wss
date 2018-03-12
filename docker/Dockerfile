FROM ubuntu
RUN apt-get update
RUN apt-get install -y curl apt-transport-https
RUN apt-get install -y openssh-server supervisor
COPY supervisor_landfw.conf /etc/supervisor/conf.d/
COPY bins/* /root/
CMD ["/root/init.sh"]
