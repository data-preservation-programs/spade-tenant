FROM golang:1.20
WORKDIR /app
COPY . .
RUN go build -o spade-tenant .
RUN go build -o spFilteringJob ./jobs
EXPOSE 8080 
CMD ["/app/spade-tenant"]