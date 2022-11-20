FROM perl:latest
RUN apt update
RUN apt install -y sudo
RUN cpanm --notest DateTime Digest::MD5 Image::ExifTool Capture::Tiny
