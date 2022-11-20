FROM perl:latest
RUN apt update
RUN apt install -y sudo
COPY . /app
WORKDIR /app
RUN cpanm --notest --installdeps .
