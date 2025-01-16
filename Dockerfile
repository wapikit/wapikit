FROM alpine:latest

# Install dependencies
RUN apk --no-cache add ca-certificates tzdata shadow su-exec

# Set the working directory
WORKDIR /wapikit

# Copy only the necessary files

# IMPORTANT: the wapikit binary must be the compatible with the target architecture which is linux/amd64
COPY wapikit .
COPY config.toml.sample config.toml

# Expose the application port
EXPOSE 8000

# Define the command to run the application
CMD ["./wapikit"]