# Start fresh from a smaller image
FROM gcr.io/distroless/static

# Copy the Pre-built binary file
COPY ./out/app /app/binary

# Expose the container to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/binary"]

