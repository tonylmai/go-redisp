# Use a base OS as a parent image
FROM centos:7

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
ADD ./go-redisp /app

# Make port 80 available to the world outside this container
EXPOSE 9997

# Define environment variable
ENV NAME RedisP

# Run the binary when the container launches
CMD /app/go-redisp