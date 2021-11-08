FROM golang:latest
WORKDIR /app
ADD . /app
RUN apt-get update -y \
&& apt-get install -y libvips-dev
# && apt-get install -y build-essential curl file git

# RUN /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" 

# ENV PATH="/home/linuxbrew/.linuxbrew/bin:$PATH"

# RUN brew install vips

ENV CGO_CFLAGS_ALLOW=-Xpreprocessor

RUN cd /app && go build
