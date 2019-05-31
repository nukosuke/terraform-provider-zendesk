## -*- docker-image-name: "terraform-zendesk" -*-
FROM hashicorp/terraform:full AS builder

ENV VERSION=master
ENV GO111MODULE="on"

WORKDIR /
ADD https://github.com/nukosuke/terraform-provider-zendesk/archive/${VERSION}.zip ./
RUN unzip ${VERSION}.zip && \
  mv terraform-provider-zendesk-${VERSION} terraform-provider-zendesk && \
  cd terraform-provider-zendesk

WORKDIR /terraform-provider-zendesk
RUN go mod tidy
RUN go mod download
RUN go build .

# dist
FROM hashicorp/terraform:light
WORKDIR /terraform
COPY --from=builder /terraform-provider-zendesk/terraform-provider-zendesk /bin/
