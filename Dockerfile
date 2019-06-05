## -*- docker-image-name: "terraform-provider-zendesk" -*-
FROM hashicorp/terraform:full AS builder

ENV GO111MODULE="on"
WORKDIR /terraform-provider-zendesk

# module cache layer
COPY go.mod go.sum /terraform-provider-zendesk/
RUN go mod tidy
RUN go mod download

# source cache layer
COPY . .
RUN go build .

# dist
FROM hashicorp/terraform:light
WORKDIR /terraform
COPY --from=builder /terraform-provider-zendesk/terraform-provider-zendesk /bin/
