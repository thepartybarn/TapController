FROM balenalib/raspberrypi3-ubuntu

WORKDIR /app

COPY main .

CMD ./main
