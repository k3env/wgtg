FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o wgtg main.go

FROM scratch
WORKDIR /app
COPY --from=builder /build/wgtg /app/wgtg
CMD ["/app/wgtg"]