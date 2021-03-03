FROM gcr.io/distroless/static:nonroot
ARG APP_NAME
COPY /${APP_NAME} /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
