FROM balenalib/raspberrypi3-ubuntu

WORKDIR /app

COPY main .

COPY web ./web

CMD ./main
